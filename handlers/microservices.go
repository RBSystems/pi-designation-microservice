package handlers

import (
	"fmt"
	"log"
	"net/http"

	ac "github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

const MICRO = "microservice_definitions"
const MMAP = "microservice_mappings"
const COLUMN_NAME = "yaml"
const MID = "microservice_id"

func AddMicroserviceDefinition(context echo.Context) error {

	log.Printf("[handlers] binding new microservice definition...")

	var microservice ac.Definition
	err := context.Bind(&microservice)
	if err != nil {
		msg := fmt.Sprintf("unable to JSON to struct", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	err = ac.AddDefinition(MICRO, &microservice)
	if err != nil {
		msg := fmt.Sprintf("unable to add microservice %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	log.Printf("%s", color.HiGreenString("[handlers] successuflly added new microservice: %s", microservice.Name))

	return context.JSON(http.StatusOK, microservice)
}

func EditMicroserviceDefinition(context echo.Context) error {

	log.Printf("[handlers] binding microservice definition...")

	var microservice ac.Definition
	err := context.Bind(&microservice)
	if err != nil {
		msg := fmt.Sprintf("unable to JSON to struct", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	log.Printf("[handlers] editing microservice definition...")

	err = ac.EditDefinition(MICRO, &microservice)
	if err != nil {
		msg := fmt.Sprintf("unable to add microservice %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	log.Printf("%s", color.HiGreenString("[handlers] successuflly added new microservice: %s", microservice.Name))

	return context.JSON(http.StatusOK, microservice)
}

func AddMicroserviceMappings(context echo.Context) error {

	log.Printf("[handlers] binding new microservice mapppings...")

	var mappings ac.Batch
	err := context.Bind(&mappings)
	if err != nil {
		msg := fmt.Sprintf("unable to bind JSON to struct: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	lastInserted, err := ac.AddMappings(MMAP, COLUMN_NAME, MID, &mappings)
	if err != nil {
		msg := fmt.Sprintf("variables not added: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	entries, err := ac.GetMicroserviceMappings(lastInserted)
	if err != nil {
		msg := fmt.Sprintf("new entries not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusInternalServerError, msg)
	}

	return context.JSON(http.StatusOK, entries)
}
