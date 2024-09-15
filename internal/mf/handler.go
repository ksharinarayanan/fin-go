package mf

import (
	"context"
	"fund-manager/internal/db"
	"fund-manager/internal/utils"
	mutual_fund "fund-manager/models/mutual_fund/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// returns all the data related to MF investments
func BaseRouteHandler(c echo.Context) error {
	mfInvestmentsApiResponse := getAllInvestments()

	return c.JSON(http.StatusOK, mfInvestmentsApiResponse)
}

func GetMfInvestmentHandler(c echo.Context) error {
	schemeIdParam := c.Param("schemeId")
	schemeId, err := strconv.Atoi(schemeIdParam)
	utils.CheckAndLogError(err, "")

	response := getInvestmentsBySchemeId(schemeId)

	return c.JSON(http.StatusOK, response)
}

func getInvestmentsBySchemeId(schemeId int) InvestmentsBySchemeIdResponse {
	ctx := context.Background()
	var response InvestmentsBySchemeIdResponse

	mfQueries := mutual_fund.New(db.DB_CONN)
	mfScheme, err := mfQueries.ListMFSchemeById(ctx, int32(schemeId))
	utils.CheckAndLogError(err, "")

	response.SchemeId = int(mfScheme.ID)
	response.SchemeName = mfScheme.SchemeName.String

	navData, err := mfQueries.ListMFNavDataBySchemeId(context.Background(), int32(schemeId))
	utils.CheckAndLogError(err, "")

	response.CurrentNav, err = utils.PgNumericToFloat64(navData[0].Nav)
	utils.CheckAndLogError(err, "")

	response.PreviousDayNav, err = utils.PgNumericToFloat64(navData[1].Nav)
	utils.CheckAndLogError(err, "")

	investments, err := mfQueries.ListMFInvestmentsBySchemeId(ctx, utils.IntToPgInt4(schemeId))
	utils.CheckAndLogError(err, "")

	response.Investments = constructInvestmentsForSchemeId(response, investments)

	return response
}

func constructInvestmentsForSchemeId(response InvestmentsBySchemeIdResponse, investments []mutual_fund.MfInvestment) []InvestmentsForSchemeId {
	var datumForSchemeId []InvestmentsForSchemeId

	for _, investment := range investments {
		var (
			investmentsForSchemeId InvestmentsForSchemeId
			err                    error
		)

		investmentsForSchemeId.Units, err = utils.PgNumericToFloat64(investment.Units)
		utils.CheckAndLogError(err, "")

		investmentsForSchemeId.InvestedAt = investment.InvestedAt.Time
		investmentsForSchemeId.InvestedNav, err = utils.PgNumericToFloat64(investment.Nav)
		utils.CheckAndLogError(err, "")

		// derived values
		investmentsForSchemeId.CurrentValue = utils.RoundFloat64(response.CurrentNav*investmentsForSchemeId.Units, NET_VALUE_ROUNDING_FACTOR)
		investmentsForSchemeId.InvestedValue = utils.RoundFloat64(investmentsForSchemeId.InvestedNav*investmentsForSchemeId.Units, NET_VALUE_ROUNDING_FACTOR)
		investmentsForSchemeId.PreviousDayValue = utils.RoundFloat64(response.PreviousDayNav*investmentsForSchemeId.Units, NET_VALUE_ROUNDING_FACTOR)

		investmentsForSchemeId.NetProfitLoss, investmentsForSchemeId.NetProfitLossPercentage = calculateProfitLoss(investmentsForSchemeId.InvestedValue, investmentsForSchemeId.CurrentValue)
		investmentsForSchemeId.DayProfitLoss, investmentsForSchemeId.DayProfitLossPercentage = calculateProfitLoss(investmentsForSchemeId.PreviousDayValue, investmentsForSchemeId.CurrentValue)

		datumForSchemeId = append(datumForSchemeId, investmentsForSchemeId)
	}

	return datumForSchemeId
}

func getAllInvestments() InvestmentsResponse {
	ctx := context.Background()

	mfQueries := mutual_fund.New(db.DB_CONN)
	schemeIds, err := mfQueries.ListDistinctMfInvestmentSchemeIds(ctx)
	utils.CheckAndLogError(err, "")

	var investmentsResponse InvestmentsResponse

	for _, schemeId := range schemeIds {
		investmentsForSchemeIdResponse := getInvestmentsBySchemeId(int(schemeId.Int32))
		investmentsResponse.Investments = append(investmentsResponse.Investments, investmentsForSchemeIdResponse)
	}

	return investmentsResponse
}

// end

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
