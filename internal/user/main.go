package main

import (
	"flag"
	"gomicroservices/internal/user/rest"
	"gomicroservices/internal/util"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	flag.Parse()
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(util.LoggerWithRequestID)

	// e.GET("api/v1/users", rest.GetUsers)
	e.GET("api/v1/users/:id", rest.GetUser)

	// Start server
	e.Logger.Fatal(e.Start(":9091"))
}
