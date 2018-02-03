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
	secure.POST("/designations/definitions", handlers.AddDesignationDefinition)
	secure.POST("/classes/definitions", handlers.AddClassDefinition)
	secure.POST("/variables/definitions", handlers.AddVariableDefinition)
	secure.POST("/microservices/definitions", handlers.AddMicroserviceDefinition)

	//add mapping
	secure.POST("/variables/mappings/multiple", handlers.AddVariableMappings)
	secure.POST("/microservices/mappings/multiple", handlers.AddMicroserviceMappings)
	secure.POST("/variables/mappings/single", handlers.AddVariableMapping)
	secure.POST("/microservices/mappings/classes/:class/designations/:designation/microservices/:microservice", handlers.AddMicroserviceMapping)

	//edit definition
	secure.PUT("/designations/definitions", handlers.EditDesignationDefinition)
	secure.PUT("/classes/definitions", handlers.EditClassDefinition)
	secure.PUT("/variables/definitions", handlers.EditVariableDefinition)
	secure.PUT("/microservices/definitions", handlers.EditMicroserviceDefinition)

	//edit mapping
	secure.PUT("/variables/mappings/single", handlers.EditVariableMapping)
	secure.PUT("/microservices/mappings/classes/:class/designations/:designation/microservices/:microservice/:mapping", handlers.EditMicroserviceMapping)

	//get definition
	secure.GET("classes/definitions/all", handlers.GetAllClassDefinitions)
	secure.GET("classes/definitions/single/:id", handlers.GetClassDefinitionById)
	secure.GET("designations/definitions/all", handlers.GetAllDesignationDefinitions)
	secure.GET("designations/definitions/single/:id", handlers.GetDesignationDefinitionById)
	secure.GET("variables/definitions/all", handlers.GetAllVariableDefinitions)
	secure.GET("variables/definitions/single/:id", handlers.GetVariableDefinitionById)
	secure.GET("microservices/definitions/all", handlers.GetAllMicroserviceDefinitions)
	secure.GET("microservices/definitions/single/:id", handlers.GetMicroserviceDefinitionById)

	//get mapping
	secure.GET("variables/mappings/all", handlers.GetAllVariableMappings)
	secure.GET("variables/mappings/single/:id", handlers.GetVariableMappingById)
	secure.GET("microservices/mappings/all", handlers.GetAllMicroserviceMappings)
	secure.GET("microservices/mappings/single/:id", handlers.GetMicroserviceMappingById)

	//delete definition
	secure.DELETE("/classes/definitions/:id", handlers.DeleteClassDefinition)
	secure.DELETE("/designations/definitions/:id", handlers.DeleteDesignationDefinition)
	secure.DELETE("/variables/definitions/:id", handlers.DeleteVariableDefinition)
	secure.DELETE("/microservices/definitions/:id", handlers.DeleteMicroserviceDefinition)

	//delete mapping
	secure.DELETE("/variables/mappings/:id", handlers.DeleteVariableMapping)
	secure.DELETE("/microservices/mappings/:id", handlers.DeleteMicroserviceMapping)

	//where the magic happens
	secure.GET("/configurations/designations/:class/:designation/variables", handlers.GetVariablesByDesignationAndClass)
	secure.GET("/configurations/designations/:class/:designation/docker-compose", handlers.GetDockerComposeByDesignationAndClass)
	secure.GET("/configurations/rooms/:room/roles/:role", handlers.GetDockerComposeByRoomAndRole)

	server := http.Server{
		Addr:           PORT,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
