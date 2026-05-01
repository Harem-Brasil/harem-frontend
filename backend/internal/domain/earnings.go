package domain

// CreatorEarningsSummaryResponse sumário agregado (sem lista de compradores).
type CreatorEarningsSummaryResponse struct {
	PeriodFrom                    string                           `json:"period_from"`
	PeriodToExclusive             string                           `json:"period_to_exclusive"`
	PlatformCommissionBasisPoints int                              `json:"platform_commission_basis_points"`
	Summaries                     []CreatorEarningsCurrencySummary `json:"summaries"`
}

// CreatorEarningsCurrencySummary totais por moeda ISO 4217.
type CreatorEarningsCurrencySummary struct {
	Currency         string `json:"currency"`
	GrossCents       int64  `json:"gross_cents"`
	PlatformFeeCents int64  `json:"platform_fee_cents"`
	NetCents         int64  `json:"net_cents"`
	PaidOrdersCount  int64  `json:"paid_orders_count"`
}
