package server

import (
	"fund-manager/internal/db"
	"fund-manager/internal/jobs"
	"fund-manager/internal/mf"
	"fund-manager/internal/utils"

	"github.com/labstack/echo/v4"
)

func StartServer() {
	db.InitializeDatabase()
	defer db.DB_CONN.Close()

	// TODO: convert this to cron job
	go jobs.UpdateMfNavData()

	e := echo.New()

	mf.RegisterRoutes(e)

	err := e.Start(":8080")
	utils.CheckAndLogError(err, "")
}
