package mf

import (
	"context"
	"fund-manager/db"
	"fund-manager/internal/utils"
	mutual_fund "fund-manager/models/mutual_fund/model"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type MFInvestment struct {
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
	MutualFunds []MFInvestment `json:"mutual_funds"`
}

// returns all the data related to MF investments
func BaseRouteHandler(c echo.Context) error {
	mfInvestmentsApiResponse := listAllMfInvestments()

	return c.JSON(http.StatusOK, mfInvestmentsApiResponse)
}

// return the JSON response to be returned
func listAllMfInvestments() MFInvestmentsApiResponse {
	mfQueries := mutual_fund.New(db.DB_CONN)
	mfInvestments, err := mfQueries.ListMFInvestments(context.Background())
	utils.CheckAndLogError(err, "")

	var mfInvestmentsApiResponse MFInvestmentsApiResponse

	for i := range mfInvestments {
		mfInvestment := mfInvestments[i]
		mfInvestmentResponse := constructMfInvestmentResponse(mfInvestment, mfQueries)

		mfInvestmentsApiResponse.MutualFunds = append(mfInvestmentsApiResponse.MutualFunds, mfInvestmentResponse)
	}

	return mfInvestmentsApiResponse
}

func constructMfInvestmentResponse(mfInvestment mutual_fund.MfInvestment, mfQueries *mutual_fund.Queries) MFInvestment {
	schemeId := int(mfInvestment.SchemeID.Int32)
	var mfInvestmentResponse MFInvestment
	mfInvestmentResponse.SchemeId = schemeId

	mfScheme, err := mfQueries.ListMFSchemeById(context.Background(), int32(schemeId))
	utils.CheckAndLogError(err, "")

	mfInvestmentResponse.SchemeName = mfScheme.SchemeName.String

	mfInvestmentResponse.Units, err = utils.PgNumericToFloat64(mfInvestment.Units)
	utils.CheckAndLogError(err, "")

	mfInvestmentResponse.InvestedAt = mfInvestment.InvestedAt.Time

	// get current navs from mf_nav_date table for the previous two days
	// this data is sorted by nav_date, so the first element is the previous day
	navDatas, err := mfQueries.ListMFNavDataBySchemeId(context.Background(), int32(schemeId))
	utils.CheckAndLogError(err, "")

	mfInvestmentResponse.CurrentNav, err = utils.PgNumericToFloat64(navDatas[0].Nav)
	utils.CheckAndLogError(err, "")
	mfInvestmentResponse.CurrentValue = utils.RoundFloat64(mfInvestmentResponse.CurrentNav*mfInvestmentResponse.Units, NET_VALUE_ROUNDING_FACTOR)

	mfInvestmentResponse.InvestedNav, err = utils.PgNumericToFloat64(mfInvestment.Nav)
	utils.CheckAndLogError(err, "")
	mfInvestmentResponse.InvestedValue = utils.RoundFloat64(mfInvestmentResponse.InvestedNav*mfInvestmentResponse.Units, NET_VALUE_ROUNDING_FACTOR)

	mfInvestmentResponse.PreviousDayNav, err = utils.PgNumericToFloat64(navDatas[1].Nav)
	utils.CheckAndLogError(err, "")
	mfInvestmentResponse.PreviousDayValue = utils.RoundFloat64(mfInvestmentResponse.PreviousDayNav*mfInvestmentResponse.Units, NET_VALUE_ROUNDING_FACTOR)

	mfInvestmentResponse.NetProfitLoss, mfInvestmentResponse.NetProfitLossPercentage = calculateProfitLoss(mfInvestmentResponse.InvestedValue, mfInvestmentResponse.CurrentValue)
	mfInvestmentResponse.DayProfitLoss, mfInvestmentResponse.DayProfitLossPercentage = calculateProfitLoss(mfInvestmentResponse.PreviousDayValue, mfInvestmentResponse.CurrentValue)

	return mfInvestmentResponse
}

// returns the P/L actual value, percentage rounded to 2 decimal places
func calculateProfitLoss(investedValue float64, currentValue float64) (float64, float64) {
	diff := currentValue - investedValue
	return utils.RoundFloat64(diff, PROFIT_LOSS_ROUNDING_FACTOR), utils.RoundFloat64(diff/investedValue*100, PROFIT_LOSS_ROUNDING_FACTOR)
}
