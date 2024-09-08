package utils

import (
	"log"
	"math"
	"time"
)

func RoundFloat64(value float64, noOfDecimalPlaces float64) float64 {
	numberToMultiplyAndDivide := math.Pow(10, noOfDecimalPlaces)
	return math.Round(value*numberToMultiplyAndDivide) / numberToMultiplyAndDivide
}

func CheckAndLogError(err error, msg string) {
	if err == nil {
		return
	}
	var errorMessage string
	if len(msg) != 0 {
		errorMessage += msg
		errorMessage += " Cause: "
	}

	errorMessage += err.Error()

	log.Fatalln(errorMessage)
}

func RemoveTimeFromDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
