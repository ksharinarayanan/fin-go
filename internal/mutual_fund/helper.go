package mutual_fund

import (
	"context"
	externalapi "fin-go/internal/mutual_fund/external_api"
	"fin-go/internal/utils"
	mutual_fund "fin-go/models/mutual_fund/model"
	"fmt"
	"strconv"
	"time"
)

// Populate's yesterday's and today's NAV in the mf_nav_data
func populateData(mfQueries *mutual_fund.Queries, mfInvestment mutual_fund.MfInvestment) error {
	todayMfNavData, yesterdayMfNavData, err := getMfNavData(mfInvestment)
	if err != nil {
		return err
	}

	err = mfQueries.AddMFNavData(context.Background(), yesterdayMfNavData)
	if err != nil {
		return err
	}

	err = mfQueries.AddMFNavData(context.Background(), todayMfNavData)
	if err != nil {
		return err
	}

	return nil
}

func getMfNavData(mfInvestment mutual_fund.MfInvestment) (mutual_fund.AddMFNavDataParams, mutual_fund.AddMFNavDataParams, error) {
	schemeId := int(mfInvestment.SchemeID.Int32)
	mfResponse, err := externalapi.MakeApiRequest(schemeId, false)

	if err != nil {
		return mutual_fund.AddMFNavDataParams{}, mutual_fund.AddMFNavDataParams{}, fmt.Errorf("error getting MF NAV data: %v", err.Error())
	}

	navDatas := mfResponse.Data

	currentDay, err := time.Parse("02-01-2006", navDatas[0].Date)
	if err != nil {
		return mutual_fund.AddMFNavDataParams{}, mutual_fund.AddMFNavDataParams{}, fmt.Errorf("error parsing MF NAV data: %v", err.Error())
	}

	previousDay, err := time.Parse("02-01-2006", navDatas[1].Date)
	if err != nil {
		return mutual_fund.AddMFNavDataParams{}, mutual_fund.AddMFNavDataParams{}, fmt.Errorf("error parsing MF NAV data: %v", err.Error())
	}

	checkIfApiIsStale(currentDay, previousDay)

	todayNavPrice, err := strconv.ParseFloat(navDatas[0].Nav, 64)
	if err != nil {
		return mutual_fund.AddMFNavDataParams{}, mutual_fund.AddMFNavDataParams{}, fmt.Errorf("error parsing MF NAV data: %v", err.Error())
	}

	yesterdayNavPrice, err := strconv.ParseFloat(navDatas[1].Nav, 64)
	if err != nil {
		return mutual_fund.AddMFNavDataParams{}, mutual_fund.AddMFNavDataParams{}, fmt.Errorf("error parsing MF NAV data: %v", err.Error())
	}

	return constructAddMFNavDataParams(schemeId, currentDay, todayNavPrice),
		constructAddMFNavDataParams(schemeId, previousDay, yesterdayNavPrice),
		nil
}

// check if the gap between yesterday and today is within one week and today is within a week of current date
// to ensure that the API is not returning stale data
func checkIfApiIsStale(currentDay time.Time, previousDay time.Time) error {
	isStale := isApartForMoreThanAWeek(currentDay, previousDay)
	if isStale {
		return fmt.Errorf("current day %v and previous day %v are more than one week apart", currentDay, previousDay)
	}

	now := time.Now()
	isStale = isApartForMoreThanAWeek(currentDay, now)
	if isStale {
		return fmt.Errorf("current day %v and now %v are more than one week apart", currentDay, now)
	}
	return nil
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
