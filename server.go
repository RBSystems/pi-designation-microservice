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

	//add definition
	secure.POST("/designations/definitions/add", handlers.AddDesignationDefinition)
	secure.POST("/classes/definitions/add", handlers.AddClassDefinition)
	secure.POST("/variables/definitions/add", handlers.AddVariableDefinition)
	secure.POST("/microservices/definitions/add", handlers.AddMicroserviceDefinition)

	//add mapping
	secure.POST("/variables/mappings/add/multiple", handlers.AddVariableMappings)
	secure.POST("/microservices/mappings/add/multiple", handlers.AddMicroserviceMappings)
	secure.POST("/variables/mappings/add/single", handlers.AddVariableMapping)
	secure.POST("/microservices/mappings/add/single", handlers.AddMicroserviceMapping)

	//edit definition
	secure.PUT("/designations/definitions/edit", handlers.EditDesignationDefinition)
	secure.PUT("/classes/definitions/edit", handlers.EditClassDefinition)
	secure.PUT("/variables/definitions/edit", handlers.EditVariableDefinition)
	secure.PUT("/microservices/definitions/edit", handlers.EditMicroserviceDefinition)

	//edit mapping
	secure.PUT("/variables/mappings/edit/single", handlers.EditVariableMapping)
	secure.PUT("/microservices/mappings/add/single", handlers.EditMicroserviceMapping)

	server := http.Server{
		Addr:           PORT,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
