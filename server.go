package main

import (
	"log"
	"net/http"

	"github.com/byuoitav/authmiddleware"
	"github.com/byuoitav/pi-designation-microservice/handlers"
	"github.com/fatih/color"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const PORT = ":5001"

func main() {

	log.Printf("%s", color.HiGreenString("Starting PI designation microservice..."))

	router := echo.New()
	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.CORS())

	secure := router.Group("", echo.WrapMiddleware(authmiddleware.Authenticate))

	//get info
	secure.GET("/rooms/:room/env", handlers.GetEnvironmentVariables)
	secure.GET("/rooms/:room/uiconfig", handlers.GetUiConfig)
	secure.GET("/variables/get/:key/:designation", handlers.GetVariable)
	secure.GET("/variables/get/:designation", handlers.GetVarsByDesignation)
	secure.GET("/variables/get/all", handlers.GetAllVariables)

	//add info
	secure.POST("/rooms/add/:room/:designation", handlers.AddNewRoom)
	secure.POST("/variables/add", handlers.AddVariable)
	secure.POST("/designations/add/:definition", handlers.AddDesignation)

	//edit info
	secure.PUT("/variables/edit/:key", handlers.EditVariable)

	//delete info
	secure.DELETE("/variables/delete/:key/:designation", handlers.DeleteVariable)
	secure.DELETE("/designations/delete/:designation", handlers.DeleteDesignation)

	server := http.Server{
		Addr:           PORT,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
