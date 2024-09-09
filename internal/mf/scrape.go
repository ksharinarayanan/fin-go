package mf

import (
	"context"
	"encoding/json"
	"fmt"
	"fund-manager/db"
	"fund-manager/internal/utils"
	mutual_fund "fund-manager/models/mutual_fund/model"
	"log"
	"net/http"
	"sync"
	"time"
)

type SchemeDataResponse struct {
	SchemeCode int    `json:"schemeCode"`
	SchemeName string `json:"schemeName"`
}

// this function is just to populate mf_schemes table
// this is only a one time activity
func ScrapeSchemeMetaData() {
	response, err := http.Get("https://api.mfapi.in/mf")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer response.Body.Close()

	var schemeDatas []SchemeDataResponse
	err = json.NewDecoder(response.Body).Decode(&schemeDatas)
	if err != nil {
		fmt.Println(err)
		return
	}

	mfQueries := mutual_fund.New(db.DB_CONN)

	var wg sync.WaitGroup
	var mutex sync.Mutex
	wg.Add(len(schemeDatas))

	counter := 0

	start := time.Now()

	for i := range schemeDatas {
		schemeData := schemeDatas[i]
		go func() {
			err := mfQueries.AddMFScheme(context.Background(), mutual_fund.AddMFSchemeParams{
				ID:         int32(schemeData.SchemeCode),
				SchemeName: utils.StringToPgText(schemeData.SchemeName),
			})

			if err != nil {
				//fmt.Printf("Error inserting scheme %v: %v\n", schemeData.SchemeCode, err.Error())
			}

			mutex.Lock()
			counter++
			if counter%100 == 0 {
				log.Printf("Processed %v records\n", counter)
			}
			mutex.Unlock()

			wg.Done()
		}()
	}

	wg.Wait()

	elapsed := time.Since(start)

	log.Printf("Processed in %v\n", elapsed)
}
