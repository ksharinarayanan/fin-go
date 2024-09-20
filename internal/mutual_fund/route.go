package mutual_fund

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterRoutes(e *echo.Echo) {
	mfRouteGroup := e.Group("/api/mf")

	mfRouteGroup.Use(middleware.CORS())

	mfRouteGroup.GET("/investments", allInvestmentsHandler)
	mfRouteGroup.POST("/investment/add", addInvestmentHandler)

	mfRouteGroup.GET("/schemes/:schemeId", getMfSchemeHandler)
}
