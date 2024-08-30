package main

import (
	"fmt"
	"log"
	"os"

	"h8-p2-finalproj-app/auth"
	"h8-p2-finalproj-app/config"
	"h8-p2-finalproj-app/service"

	// _ "h8-p2-finalproj-app/docs"
	"h8-p2-finalproj-app/handler"
	"h8-p2-finalproj-app/util"

	"github.com/golang-jwt/jwt/v5"
	echoSwagger "github.com/swaggo/echo-swagger"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

//	@title			H8 P2 Final Project App
//	@version		1.0
//	@description	Hacktiv8 Phase 2 Final Project

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/

//	@securitydefinitions.basic	BasicAuth
//	@tokenUrl					https://localhost:8080/user/login
//	@scope.read					Grants read access
//	@scope.write				Grants write access

func main() {
	db := config.CreateDBInstance()

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.HTTPErrorHandler = util.ErrorHandler

	// jwt middleware
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JwtAppClaims)
		},
		SigningKey: []byte(os.Getenv("JWT_KEY")),
	}
	jwtAuth := echojwt.WithConfig(config)

	// users
	user := handler.NewUserHandler(db)
	e.POST("/users/register", user.HandleRegisterUser)
	e.POST("/users/login", user.HandleLoginUser)

	// cars
	car := handler.NewCarHandler(service.NewCarService(db))
	cars := e.Group("/cars")
	cars.GET("", car.HandleGetCars)

	// rentals
	rental := handler.NewRentalHandler(db, service.NewCarService(db), service.NewInvoiceService())
	rentals := e.Group("/rentals")
	rentals.Use(jwtAuth)
	rentals.POST("", rental.HandlePostRentals)
	rentals.GET("", rental.HandleGetRentals)

	// payments, for call backs by xendit
	payment := handler.NewPaymentHandler(db)
	payments := e.Group("/payments")
	payments.POST("/callback", payment.HandlePaymentSuccess)

	// swagger docs
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// start server
	log.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
