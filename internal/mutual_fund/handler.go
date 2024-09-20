package mutual_fund

import (
	"context"
	"fin-go/internal/db"
	externalapi "fin-go/internal/mutual_fund/external_api"
	"fin-go/internal/utils"
	mutual_fund "fin-go/models/mutual_fund/model"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

// returns all the data related to MF investments
func baseRouteHandler(c echo.Context) error {
	mfInvestmentsApiResponse, err := getAllInvestments()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, utils.InternalServerResponse)
	}

	return c.JSON(http.StatusOK, mfInvestmentsApiResponse)
}

func getMfSchemeHandler(c echo.Context) error {
	schemeIdParam := c.Param("schemeId")
	schemeId, err := strconv.Atoi(schemeIdParam)
	if err != nil {
		if err == ErrInvalidSchemeId || err == ErrNoInvestmentsForSchemeId {
			return c.JSON(http.StatusBadRequest, utils.Response{
				Message: err.Error(),
			})
		}
		return c.JSON(http.StatusInternalServerError, utils.InternalServerResponse)
	}

	schemeData, err := getSchemeData(schemeId)
	if err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusInternalServerError, utils.InternalServerResponse)
	}

	return c.JSON(http.StatusOK, schemeData)
}

func addInvestmentHandler(c echo.Context) error {
	// post request

	return nil
}

func getInvestmentsBySchemeId(schemeId int) (InvestmentsBySchemeIdResponse, error) {
	ctx := context.Background()
	var response InvestmentsBySchemeIdResponse

	mfQueries := mutual_fund.New(db.DB_CONN)
	mfScheme, err := mfQueries.ListMFSchemeById(ctx, int32(schemeId))

	if err != nil {
		if err == pgx.ErrNoRows {
			return response, ErrInvalidSchemeId
		}
		return response, err
	}

	response.SchemeId = int(mfScheme.ID)
	response.SchemeName = mfScheme.SchemeName.String

	navData, err := mfQueries.ListMFNavDataBySchemeId(context.Background(), int32(schemeId))
	if err != nil {
		return response, err
	}

	if len(navData) == 0 {
		return response, ErrNoInvestmentsForSchemeId
	}

	response.CurrentNav, err = utils.PgNumericToFloat64(navData[0].Nav)
	if err != nil {
		return response, err
	}

	response.PreviousDayNav, err = utils.PgNumericToFloat64(navData[1].Nav)
	if err != nil {
		return response, err
	}

	investments, err := mfQueries.ListMFInvestmentsBySchemeId(ctx, utils.IntToPgInt4(schemeId))
	if err != nil {
		return response, err
	}

	response.Investments, err = constructInvestmentsForScheme(response.CurrentNav, response.PreviousDayNav, investments)
	if err != nil {
		return response, err
	}

	return response, nil
}

func constructInvestmentsForScheme(currentNav float64, previousDayNav float64, investments []mutual_fund.MfInvestment) ([]InvestmentsForScheme, error) {
	var datumForSchemeId []InvestmentsForScheme

	for _, investment := range investments {
		var (
			investmentsForScheme InvestmentsForScheme
			err                  error
		)

		investmentsForScheme.Units, err = utils.PgNumericToFloat64(investment.Units)
		if err != nil {
			log.Println(err.Error())
			return datumForSchemeId, err
		}

		investmentsForScheme.InvestedAt = investment.InvestedAt.Time
		investmentsForScheme.InvestedNav, err = utils.PgNumericToFloat64(investment.Nav)
		if err != nil {
			log.Println(err.Error())
			return datumForSchemeId, err
		}

		// derived values
		investmentsForScheme.CurrentValue = utils.RoundFloat64(currentNav*investmentsForScheme.Units, NET_VALUE_ROUNDING_FACTOR)
		investmentsForScheme.InvestedValue = utils.RoundFloat64(investmentsForScheme.InvestedNav*investmentsForScheme.Units, NET_VALUE_ROUNDING_FACTOR)
		investmentsForScheme.PreviousDayValue = utils.RoundFloat64(previousDayNav*investmentsForScheme.Units, NET_VALUE_ROUNDING_FACTOR)

		investmentsForScheme.NetProfitLoss, investmentsForScheme.NetProfitLossPercentage = calculateProfitLoss(investmentsForScheme.InvestedValue, investmentsForScheme.CurrentValue)
		investmentsForScheme.DayProfitLoss, investmentsForScheme.DayProfitLossPercentage = calculateProfitLoss(investmentsForScheme.PreviousDayValue, investmentsForScheme.CurrentValue)

		datumForSchemeId = append(datumForSchemeId, investmentsForScheme)
	}

	return datumForSchemeId, nil
}

func getAllInvestments() (InvestmentsResponse, error) {
	ctx := context.Background()

	var investmentsResponse InvestmentsResponse

	mfQueries := mutual_fund.New(db.DB_CONN)
	schemeIds, err := mfQueries.ListDistinctMfInvestmentSchemeIds(ctx)
	if err != nil {
		return investmentsResponse, err
	}

	for _, schemeId := range schemeIds {
		investmentsForSchemeIdResponse, err := getInvestmentsBySchemeId(int(schemeId.Int32))
		if err != nil {
			return investmentsResponse, err
		}
		investmentsResponse.Investments = append(investmentsResponse.Investments, investmentsForSchemeIdResponse)
	}

	return investmentsResponse, nil
}

// returns the P/L actual value and P/L percentage, percentage rounded to 2 decimal places
func calculateProfitLoss(investedValue float64, currentValue float64) (float64, float64) {
	diff := currentValue - investedValue
	return utils.RoundFloat64(diff, PROFIT_LOSS_ROUNDING_FACTOR), utils.RoundFloat64(diff/investedValue*100, PROFIT_LOSS_ROUNDING_FACTOR)
}

func getSchemeData(schemeId int) (SchemeDataResponse, error) {
	response, err := externalapi.MakeApiRequest(schemeId, true)
	var schemeData SchemeDataResponse
	if err != nil {
		log.Println(err.Error())
		return schemeData, err
	}

	schemeData.SchemeCode = response.Meta.SchemeCode
	schemeData.SchemeName = response.Meta.SchemeName
	schemeData.CurrentNav, err = strconv.ParseFloat(response.Data[0].Nav, 64)
	schemeData.Date = response.Data[0].Date
	if err != nil {
		log.Println(err.Error())
		return schemeData, err
	}

	return schemeData, nil
}
