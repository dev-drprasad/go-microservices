package main

import (
	"context"
	"fmt"
	customerrest "gomicroservices/internal/customer/rest"
	orderrest "gomicroservices/internal/order/rest"
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

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		header := c.Request().Header.Get("Authorization")
		tokenStr := strings.TrimPrefix(header, "Bearer ")

		if tokenStr == "" {
			return c.JSON(http.StatusUnauthorized, util.NewAPIError(http.StatusText(http.StatusUnauthorized)))
		}

		user, err := service.AuthCheck(tokenStr)
		if err != nil {
			c.Logger().Errorf(err.Error())
			return c.JSON(http.StatusUnauthorized, util.NewAPIError(http.StatusText(http.StatusUnauthorized)))
		}

		c.Set("user", user)

		if err := next(c); err != nil {
			c.Error(err)
		}

		return nil
	}
}

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

	e.Static("/static", "/tmp")

	userhandler := userrest.New(db)
	orghandler := orgrest.New(db)
	producthandler := productrest.New(db)
	customerhandler := customerrest.New(db)
	orderhandler := orderrest.New(db)

	e.POST("/api/v1/login", userhandler.Login)

	api := e.Group("/api/v1")
	api.Use(AuthMiddleware)

	api.GET("/users", userhandler.GetUsers)
	api.GET("/users/:id", userhandler.GetUser)
	api.POST("/users", userhandler.CreateUser)

	api.POST("/organizations", orghandler.CreateOrganization)
	api.GET("/organizations", orghandler.GetOrganizations)
	api.GET("/organizations/:id/branches", orghandler.GetBranchesByOrganization)
	api.GET("/branches", orghandler.GetBranches)
	// r.POST("/branches", orghandler.CreateBranch)

	api.POST("/brands", producthandler.CreateBrand)
	api.GET("/brands/:id", producthandler.GetBrand)
	api.GET("/brands", producthandler.GetBrands)

	api.POST("/categories", producthandler.CreateCategory)
	api.GET("/categories/:id", producthandler.GetCategory)
	api.GET("/categories", producthandler.GetCategories)

	api.POST("/products", producthandler.CreateProduct)
	api.GET("/products/:id", producthandler.GetProduct)
	api.PUT("/products/:id", producthandler.UpdateProduct)
	api.GET("/products", producthandler.GetProducts)
	api.POST("/products/images/upload", producthandler.UploadImages)

	api.POST("/orders", orderhandler.PlaceOrder)
	api.GET("/orders/:id", orderhandler.GetOrder)
	api.GET("/orders", orderhandler.GetOrders)

	api.POST("/customers", customerhandler.AddCustomer)
	api.GET("/customers/:id", customerhandler.GetCustomer)
	api.PUT("/customers/:id", customerhandler.UpdateCustomer)
	api.GET("/customers", customerhandler.GetCustomers)
	api.GET("/customers/aggregations/new", customerhandler.NewCustomersCount)

	e.Logger.Fatal(e.Start(":9090"))
}
