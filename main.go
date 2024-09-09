package main

import (
	"fund-manager/db"
	"fund-manager/internal/jobs"
	"fund-manager/internal/mf"
)

func main() {

	db.InitializeDatabase()
	defer db.DB_CONN.Close()

	// TODO: convert this to cron job
	jobs.UpdateMfNavData()

	mf.ListAllMfData()

}
