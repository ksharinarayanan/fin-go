package mutual_fund

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCheckIfStaleApi_CurrentDay(t *testing.T) {
	day1 := time.Now()
	day2 := time.Now().AddDate(0, 0, -1)

	err := checkIfApiIsStale(day1, day2)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCheckIfStaleApi_PastMonth(t *testing.T) {
	day1 := time.Now().AddDate(0, 1, 0)
	day2 := time.Now().AddDate(0, 1, -1)

	err := checkIfApiIsStale(day1, day2)
	assert.ErrorContains(t, err, "now")
}

func TestCheckIfStaleApi_DataMissingBetweenDays(t *testing.T) {
	day1 := time.Now().AddDate(0, 0, -1)
	day2 := time.Now().AddDate(0, 0, -10)

	err := checkIfApiIsStale(day1, day2)
	assert.ErrorContains(t, err, "current")
}
