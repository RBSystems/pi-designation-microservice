package handlers

import (
	"fmt"
	"log"
	"net/http"

	ac "github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

const DESIGNATION_TABLE_NAME = "designation_definitions"

func AddDesignationDefinition(context echo.Context) error {

	log.Printf("[handlers] adding new desigation definition")

	var designation ac.Definition
	err := context.Bind(&designation)
	if err != nil {
		msg := fmt.Sprintf("unable to bind JSON to struct: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusInternalServerError, msg)
	}

	err = ac.AddDefinition(DESIGNATION_TABLE_NAME, &designation)
	if err != nil {
		msg := fmt.Sprintf("error adding designation: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusInternalServerError, msg)
	}

	log.Printf("%s", color.HiGreenString("[handlers] successfully added desigation: %s", designation.Name))

	return context.JSON(http.StatusOK, designation)
}

func EditDesignationDefinition(context echo.Context) error {

	log.Printf("[handlers] editing designation definition")

	var designation ac.Definition
	err := context.Bind(&designation)
	if err != nil {
		msg := fmt.Sprintf("unable to bind JSON to struct %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	err = ac.EditDefinition(DESIGNATION_TABLE_NAME, &designation)
	if err != nil {
		msg := fmt.Sprintf("entry not updated: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, designation)
}

func GetDesignationDefinitionById(context echo.Context) error {

	id := context.Param("id")

	log.Printf("[handlers] getting designation with ID: %s", id)

	var designation ac.Definition
	err := ac.GetDefinitionById(DESIGNATION_TABLE_NAME, &designation)
	if err != nil {
		msg := fmt.Sprintf("Designation definition not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, designation)
}

func GetAllDesignationDefinitions(context echo.Context) error {

	log.Printf("[handlers] fetching all designation definitions...")

	var designations []ac.Definition
	err := ac.GetAllDefinitions(DESIGNATION_TABLE_NAME, &designations)
	if err != nil {
		msg := fmt.Sprintf("Designation definitions not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, designations)
}
