package main

import (
	"context"
	"fmt"
	orgrest "gomicroservices/internal/organization/rest"
	productrest "gomicroservices/internal/product/rest"
	userrest "gomicroservices/internal/user/rest"
	"gomicroservices/internal/user/service"
	"gomicroservices/internal/util"
	"net/http"
	"strings"

	pgx "github.com/jackc/pgx/v4/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	// connstr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", user, password, host, port, dbname)

	ctx := context.Background()
	db, err := pgx.Connect(ctx, connstr)
	if err != nil {
		panic(fmt.Errorf("unable to connect to database: %v", err))
	}
	defer db.Close()

	// db, err := sql.Open("postgres", connstr)
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close()

	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }

	// Echo instance
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(util.LoggerWithRequestID)

	userhandler := userrest.New(db)
	orghandler := orgrest.New(db)
	producthandler := productrest.New(db)

	e.POST("/api/v1/login", userhandler.Login)

	r := e.Group("/api/v1")
	r.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			header := c.Request().Header.Get("Authorization")
			tokenStr := strings.TrimPrefix(header, "Bearer ")

			if tokenStr == "" {
				return c.JSON(http.StatusUnauthorized, util.Unauthorized)
			}

			user, err := service.AuthCheck(tokenStr)
			if err != nil {
				c.Logger().Errorf(err.Error())
				return c.JSON(http.StatusUnauthorized, util.Unauthorized)
			}

			c.Set("user", user)

			if err := next(c); err != nil {
				c.Error(err)
			}

			return nil
		}
	})

	r.GET("/users", userhandler.GetUsers)
	r.GET("/users/:id", userhandler.GetUser)
	r.POST("/users", userhandler.CreateUser)

	r.POST("/organizations", orghandler.CreateOrganization)
	r.GET("/organizations", orghandler.GetOrganizations)
	r.GET("/organizations/:id/branches", orghandler.GetBranchesByOrganization)
	r.GET("/branches", orghandler.GetBranches)
	// r.POST("/branches", orghandler.CreateBranch)

	r.POST("/brands", producthandler.CreateBrand)
	r.GET("/brands/:id", producthandler.GetBrand)
	r.GET("/brands", producthandler.GetBrands)

	r.POST("/categories", producthandler.CreateCategory)
	r.GET("/categories/:id", producthandler.GetCategory)
	r.GET("/categories", producthandler.GetCategories)

	r.POST("/products", producthandler.CreateProduct)
	r.GET("/products/:id", producthandler.GetProduct)
	r.GET("/products", producthandler.GetProducts)
	r.POST("/products/images/upload", producthandler.UploadImages)

	e.Logger.Fatal(e.Start(":9090"))
}
