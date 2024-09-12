package mf

import (
	"context"
	"fund-manager/db"
	"fund-manager/internal/utils"
	mutual_fund "fund-manager/models/mutual_fund/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// returns all the data related to MF investments
func BaseRouteHandler(c echo.Context) error {
	mfInvestmentsApiResponse := listAllMfInvestments()

	return c.JSON(http.StatusOK, mfInvestmentsApiResponse)
}

func GetMfInvestmentHandler(c echo.Context) error {
	schemeIdParam := c.Param("schemeId")
	schemeId, err := strconv.Atoi(schemeIdParam)
	utils.CheckAndLogError(err, "")

	ctx := context.Background()
	var response MFInvestmentBySchemeIdResponse

	mfQueries := mutual_fund.New(db.DB_CONN)
	mfScheme, err := mfQueries.ListMFSchemeById(ctx, int32(schemeId))
	utils.CheckAndLogError(err, "")

	response.SchemeId = int(mfScheme.ID)
	response.SchemeName = mfScheme.SchemeName.String

	mfNavData, err := mfQueries.ListMFNavDataBySchemeId(context.Background(), int32(schemeId))
	utils.CheckAndLogError(err, "")

	response.CurrentNav, err = utils.PgNumericToFloat64(mfNavData[0].Nav)
	utils.CheckAndLogError(err, "")

	response.PreviousDayNav, err = utils.PgNumericToFloat64(mfNavData[1].Nav)
	utils.CheckAndLogError(err, "")

	mfInvestments, err := mfQueries.ListMFInvestmentsBySchemeId(ctx, utils.IntToPgInt4(schemeId))
	utils.CheckAndLogError(err, "")

	response.Investments = constructMfDataForSchemeId(response, mfInvestments)

	return c.JSON(http.StatusOK, response)
}

func constructMfDataForSchemeId(response MFInvestmentBySchemeIdResponse, mfinvestments []mutual_fund.MfInvestment) []MFDataForSchemeId {
	var mfDatumForSchemeId []MFDataForSchemeId

	for _, mfInvestment := range mfinvestments {
		var (
			mfDataForSchemeId MFDataForSchemeId
			err               error
		)

		mfDataForSchemeId.Units, err = utils.PgNumericToFloat64(mfInvestment.Units)
		utils.CheckAndLogError(err, "")

		mfDataForSchemeId.InvestedAt = mfInvestment.InvestedAt.Time
		mfDataForSchemeId.InvestedNav, err = utils.PgNumericToFloat64(mfInvestment.Nav)
		utils.CheckAndLogError(err, "")

		// derived values
		mfDataForSchemeId.CurrentValue = utils.RoundFloat64(response.CurrentNav*mfDataForSchemeId.Units, NET_VALUE_ROUNDING_FACTOR)
		mfDataForSchemeId.InvestedValue = utils.RoundFloat64(mfDataForSchemeId.InvestedNav*mfDataForSchemeId.Units, NET_VALUE_ROUNDING_FACTOR)
		mfDataForSchemeId.PreviousDayValue = utils.RoundFloat64(response.PreviousDayNav*mfDataForSchemeId.Units, NET_VALUE_ROUNDING_FACTOR)

		mfDataForSchemeId.NetProfitLoss, mfDataForSchemeId.NetProfitLossPercentage = calculateProfitLoss(mfDataForSchemeId.InvestedValue, mfDataForSchemeId.CurrentValue)
		mfDataForSchemeId.DayProfitLoss, mfDataForSchemeId.DayProfitLossPercentage = calculateProfitLoss(mfDataForSchemeId.PreviousDayValue, mfDataForSchemeId.CurrentValue)

		mfDatumForSchemeId = append(mfDatumForSchemeId, mfDataForSchemeId)
	}

	return mfDatumForSchemeId
}

// return the JSON response to be returned
func listAllMfInvestments() MFInvestmentsApiResponse {
	mfQueries := mutual_fund.New(db.DB_CONN)
	mfInvestments, err := mfQueries.ListMFInvestments(context.Background())
	utils.CheckAndLogError(err, "")

	var mfInvestmentsApiResponse MFInvestmentsApiResponse
	mfInvestmentsApiResponse.MutualFunds = constructMfInvestmentResponse(mfInvestments, mfQueries)

	return mfInvestmentsApiResponse
}

func constructMfInvestmentResponse(mfInvestments []mutual_fund.MfInvestment, mfQueries *mutual_fund.Queries) []MFInvestmentResponse {
	var mfInvestmentResponses []MFInvestmentResponse

	for _, mfInvestment := range mfInvestments {
		var mfInvestmentResponse MFInvestmentResponse
		schemeId := int(mfInvestment.SchemeID.Int32)
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

		mfInvestmentResponses = append(mfInvestmentResponses, mfInvestmentResponse)
	}

	return mfInvestmentResponses
}

// returns the P/L actual value, percentage rounded to 2 decimal places
func calculateProfitLoss(investedValue float64, currentValue float64) (float64, float64) {
	diff := currentValue - investedValue
	return utils.RoundFloat64(diff, PROFIT_LOSS_ROUNDING_FACTOR), utils.RoundFloat64(diff/investedValue*100, PROFIT_LOSS_ROUNDING_FACTOR)
}
