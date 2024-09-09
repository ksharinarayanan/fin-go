package mf

import "github.com/labstack/echo/v4"

func RegisterRoutes(e *echo.Echo) {
	mfRouteGroup := e.Group("/api/mf")

	mfRouteGroup.GET("/", BaseRouteHandler)
}
