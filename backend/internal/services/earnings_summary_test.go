package services

import (
	"testing"
	"time"
)

func TestSplitCommissionCents_table(t *testing.T) {
	tests := []struct {
		name    string
		gross   int64
		bps     int
		wantFee int64
		wantNet int64
	}{
		{"zero gross", 0, 1500, 0, 0},
		{"15 percent", 10000, 1500, 1500, 8500},
		{"small odd cents", 99, 1500, 14, 85},
		{"no commission", 5000, 0, 0, 5000},
		{"full platform", 1000, 10000, 1000, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fee, net := splitCommissionCents(tt.gross, tt.bps)
			if fee != tt.wantFee || net != tt.wantNet {
				t.Fatalf("splitCommissionCents(%d,%d) = (%d,%d), want (%d,%d)",
					tt.gross, tt.bps, fee, net, tt.wantFee, tt.wantNet)
			}
		})
	}
}

func TestValidateEarningsQueryParams(t *testing.T) {
	if err := validateEarningsQueryParams("", ""); err != nil {
		t.Fatal(err)
	}
	if err := validateEarningsQueryParams("2026-01-01", "2026-01-31"); err != nil {
		t.Fatal(err)
	}
	if err := validateEarningsQueryParams("2026-01-01", ""); err == nil {
		t.Fatal("expected error when only to missing")
	}
	if err := validateEarningsQueryParams("", "2026-01-01"); err == nil {
		t.Fatal("expected error when only from missing")
	}
}

func TestParseEarningsPeriod_synthetic(t *testing.T) {
	now := time.Date(2026, 6, 15, 12, 30, 0, 0, time.UTC)
	from, toExcl, err := parseEarningsPeriod("", "", now)
	if err != nil {
		t.Fatal(err)
	}
	if !toExcl.Equal(now.UTC()) {
		t.Fatalf("default toExclusive = now, got %v", toExcl)
	}
	if diff := toExcl.Sub(from); diff != defaultEarningsLookback {
		t.Fatalf("default window %v, want %v", diff, defaultEarningsLookback)
	}

	from2, toExcl2, err := parseEarningsPeriod("2026-01-01", "2026-01-31", now)
	if err != nil {
		t.Fatal(err)
	}
	wantFrom := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	wantTo := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	if !from2.Equal(wantFrom) || !toExcl2.Equal(wantTo) {
		t.Fatalf("date-only window got [%v,%v), want [%v,%v)", from2, toExcl2, wantFrom, wantTo)
	}

	from3, toExcl3, err := parseEarningsPeriod("2026-04-01T00:00:00Z", "2026-05-01T00:00:00Z", now)
	if err != nil {
		t.Fatal(err)
	}
	if !from3.Equal(time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)) || !toExcl3.Equal(time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("RFC3339 window mismatch: [%v,%v)", from3, toExcl3)
	}
}

func TestEarningsRangeExceeded(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	toOK := from.Add(365 * 24 * time.Hour)
	if earningsRangeExceeded(from, toOK) {
		t.Fatal("365d span should be allowed")
	}
	toBad := from.Add(maxEarningsRange + time.Hour)
	if !earningsRangeExceeded(from, toBad) {
		t.Fatal("span beyond 366d should be rejected")
	}
}
