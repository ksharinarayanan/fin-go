package externalapi

import (
	"encoding/json"
	"log"
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

const BASE_MF_API_URL string = "https://api.mfapi.in/mf/"

func MakeApiRequest(schemeId int, latest bool) (MFResponse, error) {
	api_url := BASE_MF_API_URL + strconv.Itoa(schemeId)
	if latest {
		api_url += "/latest"
	}
	var mfResponse MFResponse

	response, err := http.Get(api_url)
	if err != nil {
		log.Println(err.Error())
		return mfResponse, err
	}

	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&mfResponse)
	if err != nil {
		return mfResponse, err
	}

	return mfResponse, nil
}
