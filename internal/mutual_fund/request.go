package mutual_fund

import "time"

type AddInvestmentRequest struct {
	SchemeID   int       `json:"scheme_id"`
	Nav        float64   `json:"nav"`
	Units      float64   `json:"units"`
	InvestedAt time.Time `json:"invested_at"`
}
