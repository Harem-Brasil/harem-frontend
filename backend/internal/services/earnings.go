package services

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/harem-brasil/backend/internal/domain"
	"github.com/harem-brasil/backend/internal/middleware"
	"github.com/harem-brasil/backend/internal/utils"
)

const (
	// MaxEarningsQueryParamLen limite de caracteres nos query params from/to (datas RFC3339 ou YYYY-MM-DD).
	MaxEarningsQueryParamLen = 48
	maxEarningsRange         = 366 * 24 * time.Hour
	defaultEarningsLookback  = 30 * 24 * time.Hour
)

// earningEligibleStatuses pedidos que contam como receita de catálogo (exclui canceled/refunded/requested).
var earningEligibleStatuses = []string{
	domain.OrderStatusPaid,
	domain.OrderStatusFulfilled,
}

// GetCreatorEarningsSummary agrega valores por moeda no período [from, toExclusive). BOLA: só creator_id = utilizador.
func (s *Services) GetCreatorEarningsSummary(ctx context.Context, user *middleware.UserClaims, fromStr, toStr string) (*domain.CreatorEarningsSummaryResponse, error) {
	if err := validateEarningsQueryParams(fromStr, toStr); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	from, toExclusive, err := parseEarningsPeriod(strings.TrimSpace(fromStr), strings.TrimSpace(toStr), now)
	if err != nil {
		return nil, err
	}

	if earningsRangeExceeded(from, toExclusive) {
		return nil, domain.Err(http.StatusBadRequest, "Date range exceeds maximum of 366 days")
	}
	if !from.Before(toExclusive) {
		return nil, domain.Err(http.StatusBadRequest, "from must be before to")
	}

	bps := s.effectivePlatformCommissionBPS()

	rows, err := s.DB.Query(ctx,
		`SELECT currency,
		        COALESCE(SUM(amount_cents), 0)::bigint,
		        COUNT(*)::bigint
		 FROM creator_orders
		 WHERE creator_id = $1::uuid
		   AND status = ANY($2::text[])
		   AND paid_at IS NOT NULL
		   AND paid_at >= $3::timestamptz
		   AND paid_at < $4::timestamptz
		 GROUP BY currency
		 ORDER BY currency ASC`,
		user.UserID, earningEligibleStatuses, from, toExclusive,
	)
	if err != nil {
		return nil, domain.Err(500, "Database error")
	}
	defer rows.Close()

	var summaries []domain.CreatorEarningsCurrencySummary
	for rows.Next() {
		var cur string
		var gross, cnt int64
		if err := rows.Scan(&cur, &gross, &cnt); err != nil {
			continue
		}
		fee, net := splitCommissionCents(gross, bps)
		summaries = append(summaries, domain.CreatorEarningsCurrencySummary{
			Currency:         cur,
			GrossCents:       gross,
			PlatformFeeCents: fee,
			NetCents:         net,
			PaidOrdersCount:  cnt,
		})
	}

	return &domain.CreatorEarningsSummaryResponse{
		PeriodFrom:                    utils.FormatRFC3339UTC(from),
		PeriodToExclusive:             utils.FormatRFC3339UTC(toExclusive),
		PlatformCommissionBasisPoints: bps,
		Summaries:                     summaries,
	}, nil
}

func validateEarningsQueryParams(fromStr, toStr string) error {
	if utf8.RuneCountInString(fromStr) > MaxEarningsQueryParamLen || utf8.RuneCountInString(toStr) > MaxEarningsQueryParamLen {
		return domain.Err(http.StatusBadRequest, "Query parameter too long")
	}
	hasFrom := strings.TrimSpace(fromStr) != ""
	hasTo := strings.TrimSpace(toStr) != ""
	if hasFrom != hasTo {
		return domain.Err(http.StatusBadRequest, "from and to must both be set or both omitted")
	}
	return nil
}

func parseEarningsPeriod(fromStr, toStr string, now time.Time) (from, toExclusive time.Time, err error) {
	now = now.UTC()
	if fromStr == "" && toStr == "" {
		toExclusive = now
		from = toExclusive.Add(-defaultEarningsLookback)
		return from, toExclusive, nil
	}

	from, err = parseEarningsInstant(fromStr, true)
	if err != nil {
		return time.Time{}, time.Time{}, domain.Err(http.StatusBadRequest, "Invalid from datetime")
	}
	toExclusive, err = parseEarningsInstant(toStr, false)
	if err != nil {
		return time.Time{}, time.Time{}, domain.Err(http.StatusBadRequest, "Invalid to datetime")
	}
	return from, toExclusive, nil
}

// parseEarningsInstant: isFrom=true → limite inferior inclusivo; isFrom=false → limite superior exclusivo (para datas só dia YYYY-MM-DD, exclusivo = 00:00 do dia seguinte em UTC).
func parseEarningsInstant(s string, isFrom bool) (time.Time, error) {
	s = strings.TrimSpace(s)
	if len(s) == 10 {
		t, err := time.ParseInLocation("2006-01-02", s, time.UTC)
		if err != nil {
			return time.Time{}, err
		}
		if isFrom {
			return t, nil
		}
		return t.AddDate(0, 0, 1), nil
	}
	for _, layout := range []string{time.RFC3339, time.RFC3339Nano} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC(), nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid datetime")
}

func earningsRangeExceeded(from, toExclusive time.Time) bool {
	return toExclusive.Sub(from) > maxEarningsRange
}

func splitCommissionCents(gross int64, platformCommissionBPS int) (fee int64, net int64) {
	if platformCommissionBPS <= 0 {
		return 0, gross
	}
	if platformCommissionBPS >= 10000 {
		return gross, 0
	}
	fee = (gross * int64(platformCommissionBPS)) / 10000
	return fee, gross - fee
}

func (s *Services) effectivePlatformCommissionBPS() int {
	bps := s.PlatformCommissionBasisPoints
	if bps <= 0 || bps > 10000 {
		return 1500
	}
	return bps
}
