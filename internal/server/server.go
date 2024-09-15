package server

import (
	"fin-go/internal/db"
	"fin-go/internal/jobs"
	"fin-go/internal/mf"
	"fin-go/internal/utils"

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
