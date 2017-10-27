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

	log.Printf("%s", color.HiGreenString("Starting room designation microservice..."))

	router := echo.New()
	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.CORS())

	secure := router.Group("", echo.WrapMiddleware(authmiddleware.Authenticate))

	secure.POST("/designations/definitions/add", handlers.AddDesignationDefinition)
	secure.POST("/classes/definitions/add", handlers.AddClassDefinition)
	secure.POST("/microservices/definitions/add", handlers.AddMicroserviceDefinition)
	secure.POST("/microservices/mappings/add", handlers.AddMicroserviceMappings)
	secure.POST("/variables/definitions/add", handlers.AddVariableDefinition)
	secure.POST("/variables/mappings/add", handlers.AddVariableMappings)

	secure.PUT("/designations/definitions/edit", handlers.EditDesignationDefinition)
	secure.PUT("/classes/definitions/edit", handlers.EditClassDefinition)

	server := http.Server{
		Addr:           PORT,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
