package mf

import (
	"context"
	"fund-manager/db"
	"fund-manager/internal/utils"
	mutual_fund "fund-manager/models/mutual_fund/model"
	"log"
	"time"
)

type MutualFundResponse struct {
	// all the values and P/L percentages are derived values
	// they can be computed in the client as well
	// TODO: is it better to offload that to client?
	SchemeId                int       `json:"scheme_id"`
	SchemeName              string    `json:"scheme_name"`
	Units                   float64   `json:"units"`
	InvestedAt              time.Time `json:"invested_at"`
	CurrentNav              float64   `json:"current_nav"`
	CurrentValue            float64   `json:"current_value"`
	InvestedNav             float64   `json:"invested_nav"`
	InvestedValue           float64   `json:"invested_value"`
	PreviousDayNav          float64   `json:"previous_day_nav"`
	PreviousDayValue        float64   `json:"previous_day_value"`
	NetProfitLossPercentage float64   `json:"net_profit_loss_percentage"`
	DayProfitLossPercentage float64   `json:"day_profit_loss_percentage"`
	NetProfitLoss           float64   `json:"net_profit_loss"`
	DayProfitLoss           float64   `json:"day_profit_loss"`
}

type ResponseObject struct {
	MutualFunds []MutualFundResponse `json:"mutual_funds"`
}

func ListAllMfData() {
	mfQueries := mutual_fund.New(db.DB_CONN)
	mfInvestments, err := mfQueries.ListMFInvestments(context.Background())
	utils.CheckAndLogError(err, "")

	for i := range mfInvestments {
		mfInvestment := mfInvestments[i]

		schemeId := int(mfInvestment.SchemeID.Int32)

		var mfResponse MutualFundResponse
		mfResponse.SchemeId = schemeId

		mfScheme, err := mfQueries.ListMFSchemeById(context.Background(), int32(schemeId))
		utils.CheckAndLogError(err, "")

		mfResponse.SchemeName = mfScheme.SchemeName.String

		mfResponse.Units, err = utils.PgNumericToFloat64(mfInvestment.Units)
		utils.CheckAndLogError(err, "")

		mfResponse.InvestedAt = mfInvestment.InvestedAt.Time

		// get current navs from mf_nav_date table for the previous two days
		// this data is sorted by nav_date, so the first element is the previous day
		navDatas, err := mfQueries.ListMFNavDataBySchemeId(context.Background(), int32(schemeId))
		utils.CheckAndLogError(err, "")

		mfResponse.CurrentNav, err = utils.PgNumericToFloat64(navDatas[0].Nav)
		utils.CheckAndLogError(err, "")
		mfResponse.CurrentValue = utils.RoundFloat64(mfResponse.CurrentNav*mfResponse.Units, NET_VALUE_ROUNDING_FACTOR)

		mfResponse.InvestedNav, err = utils.PgNumericToFloat64(mfInvestment.Nav)
		utils.CheckAndLogError(err, "")
		mfResponse.InvestedValue = utils.RoundFloat64(mfResponse.InvestedNav*mfResponse.Units, NET_VALUE_ROUNDING_FACTOR)

		mfResponse.PreviousDayNav, err = utils.PgNumericToFloat64(navDatas[1].Nav)
		utils.CheckAndLogError(err, "")
		mfResponse.PreviousDayValue = utils.RoundFloat64(mfResponse.PreviousDayNav*mfResponse.Units, NET_VALUE_ROUNDING_FACTOR)

		mfResponse.NetProfitLoss, mfResponse.NetProfitLossPercentage = calculateProfitLoss(mfResponse.InvestedValue, mfResponse.CurrentValue)
		mfResponse.DayProfitLoss, mfResponse.DayProfitLossPercentage = calculateProfitLoss(mfResponse.PreviousDayValue, mfResponse.CurrentValue)

		print(mfResponse)
	}

}

// returns the P/L actual value, percentage rounded to 2 decimal places
func calculateProfitLoss(investedValue float64, currentValue float64) (float64, float64) {
	diff := currentValue - investedValue
	return utils.RoundFloat64(diff, PROFIT_LOSS_ROUNDING_FACTOR), utils.RoundFloat64(diff/investedValue*100, PROFIT_LOSS_ROUNDING_FACTOR)
}

func print(mfResponse MutualFundResponse) {
	log.Printf("Scheme ID: %v\n", mfResponse.SchemeId)
	log.Printf("Scheme Name: %v\n", mfResponse.SchemeName)
	log.Printf("Units: %v\n", mfResponse.Units)
	log.Printf("Invested At: %v\n", mfResponse.InvestedAt)
	log.Printf("Current NAV: %v\n", mfResponse.CurrentNav)
	log.Printf("Current Value: %v\n", mfResponse.CurrentValue)
	log.Printf("Invested NAV: %v\n", mfResponse.InvestedNav)
	log.Printf("Invested Value: %v\n", mfResponse.InvestedValue)
	log.Printf("Previous Day NAV: %v\n", mfResponse.PreviousDayNav)
	log.Printf("Previous Day Value: %v\n", mfResponse.PreviousDayValue)
	log.Printf("Net Profit/Loss: %v\n", mfResponse.NetProfitLoss)
	log.Printf("Day Profit/Loss: %v\n", mfResponse.DayProfitLoss)
	log.Printf("Net Profit/Loss Percentage: %v\n", mfResponse.NetProfitLossPercentage)
	log.Printf("Day Profit/Loss Percentage: %v\n", mfResponse.DayProfitLossPercentage)
}
