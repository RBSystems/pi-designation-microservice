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
	secure.PUT("/microservices/mappings/edit/single", handlers.EditMicroserviceMapping) //TODO test this and all that follow

	//get definition
	secure.GET("class/definitions/all", handlers.GetAllClassDefinitions)
	secure.GET("class/definitions/single/:id", handlers.GetClassDefinitionById)
	secure.GET("designations/definitions/all", handlers.GetAllDesignationDefinitions)
	secure.GET("designations/definitions/single/:id", handlers.GetDesignationDefinitionById)
	secure.GET("variables/definitions/all", handlers.GetVariableDefinition)
	secure.GET("variables/definitions/single/:id", handlers.GetAllVariableDefinitions)
	secure.GET("microservices/definitions/all", handlers.GetAllMicroserviceDefinitions)
	secure.GET("microservices/definitions/single/:id", handlers.GetMicroserviceDefinitionById)

	//get mapping
	secure.GET("variables/mappings/all", handlers.GetAllVariableMappings)
	secure.GET("variables/mappings/single/:id", handlers.GetVariableMappingById)
	secure.GET("microservices/mappings/all", handlers.GetAllVariableMappings)
	secure.GET("microservices/mappings/single/:id", handlers.GetVariableMappingById)

	//get mappings

	server := http.Server{
		Addr:           PORT,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
