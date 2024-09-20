package server

import (
	"fin-go/internal/db"
	"fin-go/internal/mutual_fund"
	"fin-go/internal/utils"

	"github.com/labstack/echo/v4"
)

func StartServer() {
	db.InitializeDatabase()
	defer db.DB_CONN.Close()

	// TODO: convert this to cron job
	mutual_fund.UpdateMfNavData()

	e := echo.New()

	mutual_fund.RegisterRoutes(e)

	err := e.Start(":8080")
	utils.CheckAndLogError(err, "")
}
