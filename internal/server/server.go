package server

import (
	"fin-go/internal/db"
	"fin-go/internal/mutual_fund"
	"fin-go/internal/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func StartServer() {
	db.InitializeDatabase()
	defer db.DB_CONN.Close()

	// TODO: convert this to cron job
	err := mutual_fund.UpdateMfNavData()
	utils.CheckAndLogError(err, "Error while updating MF nav data: ")

	e := echo.New()

	e.Use(middleware.Logger())
	mutual_fund.RegisterRoutes(e)

	e.Start(":8080")
}
