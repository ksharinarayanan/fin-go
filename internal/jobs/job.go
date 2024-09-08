package jobs

import (
	"context"
	"fund-manager/db"
	"fund-manager/internal/mf"
	"fund-manager/internal/utils"
	mutual_fund "fund-manager/models/mutual_fund/model"
)

// clean existing data and populate yesterday's and today's NAV in the mf_nav_data
func UpdateMfNavData() {
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
		mfQueries.CleanupMFNavDataBySchemeId(context.Background(), mfInvestments[0].SchemeID.Int32)
		mf.PopulateData(mfQueries, mfInvestments[i])
	}

	tx.Commit(context.Background())
}
