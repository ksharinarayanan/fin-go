package mf

import (
	"context"
	"fund-manager/db"
	"fund-manager/internal/utils"
	mutual_fund "fund-manager/models/mutual_fund/model"
	"log"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
)

func Scrape() {
	// jobs.UpdateMfNavData()
	// var mfQueries *mutual_fund.Queries = mutual_fund.New(db.DB_CONN)

	// mfInvestments, err := mfQueries.ListMFInvestments(context.Background())
	// utils.CheckAndLogError(err, "")

	// for i := range mfInvestments {
	// 	processScheme(mfInvestments[i])
	// }
}

func processScheme(mfInvestment mutual_fund.MfInvestment) {
	schemeId := int(mfInvestment.SchemeID.Int32)

	log.Printf("Processing scheme ID %v\n", schemeId)

	mfResponse := apiRequest(schemeId, true)

	var navDatas []NavData = mfResponse.Data

	nav, err := utils.PgNumericToFloat64(mfInvestment.Nav)
	utils.CheckAndLogError(err, "")
	units, err := utils.PgNumericToFloat64(mfInvestment.Units)
	utils.CheckAndLogError(err, "")
	navToday, err := strconv.ParseFloat(navDatas[0].Nav, 64)
	utils.CheckAndLogError(err, "")

	investedValue := utils.RoundFloat64(nav*units, 3)
	currentValue := utils.RoundFloat64(navToday*units, 3)

	compareToYesterday(schemeId, units, currentValue)

	log.Printf("Net P/L: %v\n", compareChangeInPercentage(investedValue, currentValue))
}

func compareChangeInPercentage(investedValue float64, currentValue float64) float32 {
	differenceInValue := currentValue - investedValue

	return float32(utils.RoundFloat64((differenceInValue/investedValue)*100, 2))
}

func compareToYesterday(schemeId int, units float64, investedValue float64) {
	var mfQueries *mutual_fund.Queries = mutual_fund.New(db.DB_CONN)

	yesterday := time.Now().AddDate(0, 0, -1)
	yesterdaysNavData, err := mfQueries.ListMFNavData(context.Background(), mutual_fund.ListMFNavDataParams{
		SchemeID: int32(schemeId),
		NavDate:  utils.TimeToPgDate(yesterday),
	})

	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("No matching rows found for scheme ID %v and date %v", schemeId, yesterday)
			return
		}
		utils.CheckAndLogError(err, "")
		return
	}

	nav, err := yesterdaysNavData.Nav.Float64Value()
	utils.CheckAndLogError(err, "")

	yesterdaysValue := nav.Float64 * units

	change := compareChangeInPercentage(yesterdaysValue, investedValue)

	log.Println(change)

}
