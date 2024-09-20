package mutual_fund

import (
	"context"
	"fin-go/internal/db"
	mutual_fund "fin-go/models/mutual_fund/model"
	"log"
)

// clean existing data and populate yesterday's and today's NAV in the mf_nav_data
func UpdateMfNavData() error {
	log.Println("Updating MF nav data")
	// use transaction
	tx, err := db.DB_CONN.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	mfQueries := mutual_fund.New(db.DB_CONN).WithTx(tx)

	// get all the schemes that's being tracked
	mfInvestments, err := mfQueries.ListMFInvestments(context.Background())
	if err != nil {
		return err
	}

	anyFailure := false

	for i := range mfInvestments {
		// remove any existing data for the scheme ID
		mfQueries.CleanupMFNavDataBySchemeId(context.Background(), mfInvestments[i].SchemeID.Int32)
		err := populateData(mfQueries, mfInvestments[i])
		if err != nil {
			anyFailure = true
			log.Printf("Error while populating data for %v: %v", mfInvestments[i].SchemeID, err)
		}
	}

	tx.Commit(context.Background())
	if anyFailure {
		log.Println("Some updates failed")
	}

	log.Println("MF nav data updation job completed")
	return nil
}
