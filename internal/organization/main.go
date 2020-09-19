package main

import (
	"flag"
	"fmt"
	"gomicroservices/internal/common/genproto/organization"
	orggrpc "gomicroservices/internal/organization/grpc"
	"gomicroservices/internal/organization/rest"
	"log"
	"net"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	grpc "google.golang.org/grpc"
)

func main() {
	flag.Parse()
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// e.GET("api/v1/users", rest.GetUsers)
	e.GET("api/v1/branches/:id", rest.GetBranch)

	// Start server
	go func() {
		e.Logger.Fatal(e.Start(":9092"))
	}()
	log.Println("Initializing grpc server")
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", 8082))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		panic(fmt.Sprintf("failed to listen: %v", err))
	}

	grpcserver := grpc.NewServer()

	organizationgrpc := orggrpc.NewOrganization()
	organization.RegisterOrganizationServiceServer(grpcserver, organizationgrpc)

	if err := grpcserver.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
