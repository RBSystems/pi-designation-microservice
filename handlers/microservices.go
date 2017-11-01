package handlers

import (
	"fmt"
	"log"
	"net/http"

	ac "github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

const MICROSERVICE_DEFINITION_TABLE = "microservice_definitions"
const MICROSERVICE_DEFINITION_COLUMN = "microservice_id"
const MICROSERVICE_MAPPINGS_TABLE = "microservice_mappings"
const MICROSERVICE_COLUMN_NAME = "yaml"

func AddMicroserviceDefinition(context echo.Context) error {

	log.Printf("[handlers] binding new microservice definition...")

	var microservice ac.Definition
	err := context.Bind(&microservice)
	if err != nil {
		msg := fmt.Sprintf("unable to JSON to struct", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	err = ac.AddDefinition(MICROSERVICE_DEFINITION_TABLE, &microservice)
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

	err = ac.EditDefinition(MICROSERVICE_DEFINITION_TABLE, &microservice)
	if err != nil {
		msg := fmt.Sprintf("unable to add microservice %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	log.Printf("%s", color.HiGreenString("[handlers] successuflly added new microservice: %s", microservice.Name))

	return context.JSON(http.StatusOK, microservice)
}

func AddMicroserviceMapping(context echo.Context) error {

	log.Printf("[handlers] binding microservice mapping...")

	var mapping ac.MicroserviceMapping
	err := context.Bind(&mapping)
	if err != nil {
		msg := fmt.Sprintf("unable to bind JSON to struct: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	id, err := ac.AddMapping(
		MICROSERVICE_MAPPINGS_TABLE,
		MICROSERVICE_DEFINITION_COLUMN,
		MICROSERVICE_COLUMN_NAME,
		mapping.YAML,
		mapping.Microservice.ID,
		mapping.Class.ID,
		mapping.Designation.ID)
	if err != nil {
		msg := fmt.Sprintf("unable to add microservice mapping: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	entry, err := ac.GetMicroserviceMapping(id)
	if err != nil {
		msg := fmt.Sprintf("mapping entry not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, entry)
}

func EditMicroserviceMapping(context echo.Context) error {

	log.Printf("[handlers] binding microservice mapping...")

	var mapping ac.MicroserviceMapping
	err := context.Bind(&mapping)
	if err != nil {
		msg := fmt.Sprintf("unable to bind JSON to struct: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	err = ac.EditMapping(
		MICROSERVICE_MAPPINGS_TABLE,
		MICROSERVICE_DEFINITION_COLUMN,
		MICROSERVICE_COLUMN_NAME,
		mapping.YAML,
		mapping.Microservice.ID,
		mapping.Class.ID,
		mapping.Designation.ID,
		mapping.ID)
	if err != nil {
		msg := fmt.Sprintf("unable edit mapping: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	entry, err := ac.GetMicroserviceMapping(mapping.ID)
	if err != nil {
		msg := fmt.Sprintf("new entries not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusInternalServerError, msg)
	}

	return context.JSON(http.StatusOK, entry)
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

	lastInserted, err := ac.AddMappings(
		MICROSERVICE_MAPPINGS_TABLE,
		MICROSERVICE_DEFINITION_COLUMN,
		MICROSERVICE_COLUMN_NAME,
		&mappings)
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
