package mf

import (
	"context"
	"fmt"
	"fund-manager/internal/utils"
	mutual_fund "fund-manager/models/mutual_fund/model"
	"strconv"
	"time"
)

// Populate's yesterday's and today's NAV in the mf_nav_data
func PopulateData(mfQueries *mutual_fund.Queries, mfInvestment mutual_fund.MfInvestment) {
	todayMfNavData, yesterdayMfNavData := getMfNavData(mfInvestment)

	mfQueries.AddMFNavData(context.Background(), yesterdayMfNavData)
	mfQueries.AddMFNavData(context.Background(), todayMfNavData)
}

func getMfNavData(mfInvestment mutual_fund.MfInvestment) (mutual_fund.AddMFNavDataParams, mutual_fund.AddMFNavDataParams) {
	schemeId := int(mfInvestment.SchemeID.Int32)
	mfResponse := apiRequest(schemeId, false)

	navDatas := mfResponse.Data

	today, err := time.Parse("02-01-2006", navDatas[0].Date)
	utils.CheckAndLogError(err, "")

	yesterday, err := time.Parse("02-01-2006", navDatas[1].Date)
	utils.CheckAndLogError(err, "")

	checkIfStaleApi(today, yesterday)

	todayNavPrice, err := strconv.ParseFloat(navDatas[0].Nav, 64)
	utils.CheckAndLogError(err, "")

	yesterdayNavPrice, err := strconv.ParseFloat(navDatas[1].Nav, 64)
	utils.CheckAndLogError(err, "")

	return constructAddMFNavDataParams(schemeId, today, todayNavPrice), constructAddMFNavDataParams(schemeId, yesterday, yesterdayNavPrice)
}

// check if the gap between yesterday and today is within one week and today is within a week of current date
// to ensure that the API is not return stale data
func checkIfStaleApi(today time.Time, yesterday time.Time) {
	isStale := isApartForMoreThanAWeek(today, yesterday)
	if isStale {
		err := fmt.Errorf("today %v and yesterday %v are more than one week apart", today, yesterday)
		utils.CheckAndLogError(err, "")
	}

	now := time.Now()
	isStale = isApartForMoreThanAWeek(today, now)
	if isStale {
		err := fmt.Errorf("today %v and Now %v are more than one week apart", today, now)
		utils.CheckAndLogError(err, "")
	}
}

func isApartForMoreThanAWeek(d1 time.Time, d2 time.Time) bool {
	diff := d1.Sub(d2)
	if diff < 0 {
		diff = -diff
	}

	return diff > time.Hour*24*7
}

func constructAddMFNavDataParams(schemeId int, date time.Time, nav float64) mutual_fund.AddMFNavDataParams {
	return mutual_fund.AddMFNavDataParams{
		SchemeID: int32(schemeId),
		NavDate:  utils.TimeToPgDate(date),
		Nav:      utils.Float64ToPgNumeric(nav),
	}
}
