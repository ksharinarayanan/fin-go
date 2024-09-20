package mutual_fund

import (
	"context"
	"encoding/json"
	"fin-go/internal/db"
	"fin-go/internal/utils"
	mutual_fund "fin-go/models/mutual_fund/model"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type SchemeMetadataResponse struct {
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

	var schemeDatas []SchemeMetadataResponse
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

// function to read all the mf_schemes data and dump it in json file for
// front end to consume
func DumpJsonFile() {
	db.InitializeDatabase()
	defer db.DB_CONN.Close()

	rows, err := db.DB_CONN.Query(context.Background(), "SELECT * FROM mf_schemes")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	var mfSchemes []mutual_fund.MfScheme

	file, err := os.Create("mf_schemes.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for rows.Next() {
		var mfScheme mutual_fund.MfScheme
		err := rows.Scan(&mfScheme.ID, &mfScheme.SchemeName)
		if err != nil {
			log.Fatal(err)
		}
		mfSchemes = append(mfSchemes, mfScheme)
	}

	err = json.NewEncoder(file).Encode(mfSchemes)
	if err != nil {
		log.Fatal(err)
	}

}
