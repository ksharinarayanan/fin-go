package mf

import "time"

// contains the request and response types for /api/mf/ group

type InvestmentsBySchemeIdResponse struct {
	SchemeId       int                      `json:"scheme_id"`
	SchemeName     string                   `json:"scheme_name"`
	CurrentNav     float64                  `json:"current_nav"`
	PreviousDayNav float64                  `json:"previous_day_nav"`
	Investments    []InvestmentsForSchemeId `json:"investments"`
}

type InvestmentsForSchemeId struct {
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

// end

type MFInvestmentResponse struct {
	// all the values and P/L percentages are derived values
	// they can be computed in the client as well
	// TODO: is it better to offload that to client?
	SchemeId                int       `json:"scheme_id"`
	SchemeName              string    `json:"scheme_name"`
	Units                   float64   `json:"units"`
	InvestedAt              time.Time `json:"invested_at"`
	CurrentNav              float64   `json:"current_nav"`
	InvestedNav             float64   `json:"invested_nav"`
	PreviousDayNav          float64   `json:"previous_day_nav"`
	NetProfitLossPercentage float64   `json:"net_profit_loss_percentage"`

	// these are derived values
	CurrentValue            float64 `json:"current_value"`
	InvestedValue           float64 `json:"invested_value"`
	PreviousDayValue        float64 `json:"previous_day_value"`
	DayProfitLossPercentage float64 `json:"day_profit_loss_percentage"`
	NetProfitLoss           float64 `json:"net_profit_loss"`
	DayProfitLoss           float64 `json:"day_profit_loss"`
}

type MFInvestmentsApiResponse struct {
	MutualFunds []MFInvestmentResponse `json:"mutual_funds"`
}

type MFInvestmentsApiRequest struct {
	// type here
}
