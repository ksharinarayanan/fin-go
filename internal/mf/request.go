package mf

import (
	"encoding/json"
	"fin-go/internal/utils"
	"fmt"
	"net/http"
	"strconv"
)

type MFResponse struct {
	Meta   SchemeData `json:"meta"`
	Data   []NavData  `json:"data"`
	Status string     `json:"status"`
}

type SchemeData struct {
	FundHouse      string `json:"fund_house"`
	SchemeType     string `json:"scheme_type"`
	SchemeCategory string `json:"scheme_category"`
	SchemeCode     int    `json:"scheme_code"`
	SchemeName     string `json:"scheme_name"`
}

type NavData struct {
	Date string `json:"date"`
	Nav  string `json:"nav"`
}

func apiRequest(schemeId int, latest bool) MFResponse {
	api_url := BASE_MF_API_URL + strconv.Itoa(schemeId)
	if latest {
		api_url += "/latest"
	}

	response, err := http.Get(api_url)
	if err != nil {
		utils.CheckAndLogError(err, fmt.Sprintf("Failed to hit MF API %v: %v\n", api_url, err.Error()))
	}

	defer response.Body.Close()

	var mfResponse MFResponse
	err = json.NewDecoder(response.Body).Decode(&mfResponse)
	if err != nil {
		utils.CheckAndLogError(err, fmt.Sprintf("Error parsing the MF API %v: %v\n", api_url, err.Error()))
	}

	return mfResponse
}
