package mutual_fund

import "time"

// contains the response types for /api/mf/ group

type InvestmentsBySchemeIdResponse struct {
	SchemeId       int                    `json:"scheme_id"`
	SchemeName     string                 `json:"scheme_name"`
	CurrentNav     float64                `json:"current_nav"`
	PreviousDayNav float64                `json:"previous_day_nav"`
	Investments    []InvestmentsForScheme `json:"investments"`
}

type InvestmentsForScheme struct {
	Units       float64   `json:"units"`
	InvestedAt  time.Time `json:"invested_at"`
	InvestedNav float64   `json:"invested_nav"`

	// these are derived values
	CurrentValue            float64 `json:"current_value"`
	InvestedValue           float64 `json:"invested_value"`
	PreviousDayValue        float64 `json:"previous_day_value"`
	NetProfitLossPercentage float64 `json:"net_profit_loss_percentage"`
	DayProfitLossPercentage float64 `json:"day_profit_loss_percentage"`
	NetProfitLoss           float64 `json:"net_profit_loss"`
	DayProfitLoss           float64 `json:"day_profit_loss"`
}

type InvestmentsResponse struct {
	Investments []InvestmentsBySchemeIdResponse `json:"investments"`
}

type SchemeDataResponse struct {
	SchemeCode int     `json:"schemeCode"`
	SchemeName string  `json:"schemeName"`
	CurrentNav float64 `json:"currentNav"`
	Date       string  `json:"date"`
}
