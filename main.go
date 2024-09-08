package main

import (
	"fund-manager/db"
	"fund-manager/internal/jobs"
)

func main() {

	db.InitializeDatabase()
	defer db.DB_CONN.Close()

	// mf.Scrape()

	// TODO: convert this to cron job
	jobs.UpdateMfNavData()

}
