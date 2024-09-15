package mf

import (
	"fund-manager/internal/utils"
	mutual_fund "fund-manager/models/mutual_fund/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConstructInvestmentsBySchemeIdResponse(t *testing.T) {
	/*
		Scenario: construct investments by scheme id response

		Scheme ID: 1
		Scheme Name: test
		Current NAV: 200
		Previous day NAV: 190

		There's two investments,
		1. One 2 months back (10 units) at 100 NAV
		2. One 1 month back (10 units) at 200 NAV
	*/
	investments := []mutual_fund.MfInvestment{
		{
			ID:         1,
			SchemeID:   utils.IntToPgInt4(1),
			Nav:        utils.Float64ToPgNumeric(100),
			Units:      utils.Float64ToPgNumeric(10),
			InvestedAt: utils.TimeToPgDate(truncateDateFromTime(time.Now().AddDate(0, 2, 0))),
		},
		{
			ID:         2,
			SchemeID:   utils.IntToPgInt4(1),
			Nav:        utils.Float64ToPgNumeric(200),
			Units:      utils.Float64ToPgNumeric(10),
			InvestedAt: utils.TimeToPgDate(truncateDateFromTime(time.Now().AddDate(0, 1, 0))),
		},
	}
	currentNav := float64(200)
	previousDayNav := float64(190)

	investmentsForSchemeId, err := constructInvestmentsForScheme(currentNav, previousDayNav, investments)
	if err != nil {
		t.Fatal(err)
	}

	expectedInvestmentsForSchemeId := []InvestmentsForScheme{
		{
			Units:       float64(10),
			InvestedAt:  truncateDateFromTime(time.Now().AddDate(0, 2, 0)),
			InvestedNav: float64(100),

			CurrentValue:            float64(2000),
			InvestedValue:           float64(1000),
			PreviousDayValue:        float64(1900),
			NetProfitLossPercentage: float64(100),
			DayProfitLossPercentage: float64(5.26),
			NetProfitLoss:           float64(1000),
			DayProfitLoss:           float64(100),
		},
		{
			Units:       float64(10),
			InvestedAt:  truncateDateFromTime(time.Now().AddDate(0, 1, 0)),
			InvestedNav: float64(200),

			CurrentValue:            float64(2000),
			InvestedValue:           float64(2000),
			PreviousDayValue:        float64(1900),
			NetProfitLossPercentage: float64(0),
			DayProfitLossPercentage: float64(5.26),
			NetProfitLoss:           float64(0),
			DayProfitLoss:           float64(100),
		},
	}

	assert.Equal(t, expectedInvestmentsForSchemeId, investmentsForSchemeId)
}

func TestCalculateProfitLoss(t *testing.T) {
	investedValue, currentValue := 100, 150

	diff, profitLossPercentage := calculateProfitLoss(float64(investedValue), float64(currentValue))

	assert.Equal(t, diff, float64(50))
	assert.Equal(t, profitLossPercentage, float64(50))
}

func truncateDateFromTime(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
