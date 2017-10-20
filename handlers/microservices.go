package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

func AddMicroserviceDefinition(context echo.Context) error {

	log.Printf("[handlers] unmarshalling new microservice definition...")

	var microservice accessors.MicroserviceDefinition
	err := context.Bind(&microservice)
	if err != nil {
		msg := fmt.Sprintf("error unmarshalling struct: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	return nil
}

func AddMicroserviceMapping(context echo.Context) error {

	log.Printf("[handlers] unmarshalling new microservice mappping...")

	return nil
}
