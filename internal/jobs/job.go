package jobs

import (
	"context"
	"fin-go/internal/db"
	"fin-go/internal/mf"
	"fin-go/internal/utils"
	mutual_fund "fin-go/models/mutual_fund/model"
	"log"
)

// clean existing data and populate yesterday's and today's NAV in the mf_nav_data
func UpdateMfNavData() {
	log.Println("Updating MF nav data")
	// use transaction
	tx, err := db.DB_CONN.Begin(context.Background())
	utils.CheckAndLogError(err, "")

	defer tx.Rollback(context.Background())

	mfQueries := mutual_fund.New(db.DB_CONN).WithTx(tx)

	// get all the schemes that's being tracked
	mfInvestments, err := mfQueries.ListMFInvestments(context.Background())
	utils.CheckAndLogError(err, "")

	for i := range mfInvestments {
		// remove any existing data for the scheme ID
		mfQueries.CleanupMFNavDataBySchemeId(context.Background(), mfInvestments[i].SchemeID.Int32)
		mf.PopulateData(mfQueries, mfInvestments[i])
	}

	tx.Commit(context.Background())
	log.Println("Updated")
}
