package main

import (
	"database/sql"
	"fmt"
	userctrl "gomicroservices/internal/user/rest"
	"gomicroservices/internal/util"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = ""
	dbname   = "postgres"
)

func main() {
	connstr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connstr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// Echo instance
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(util.LoggerWithRequestID)

	userhandler := userctrl.New(db)

	// e.GET("api/v1/users", rest.GetUsers)
	e.GET("/api/v1/users/:id", userhandler.GetUser)
	e.POST("/api/v1/users", userhandler.CreateUser)

	// Start server
	e.Logger.Fatal(e.Start(":9090"))
}
