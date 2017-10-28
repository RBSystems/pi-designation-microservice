package handlers

import (
	"fmt"
	"log"
	"net/http"

	ac "github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

const VAR = "variable_definitions"
const VMAP = "variable_mappings"
const VID = "variable_id"

//relies on MySQL for most logic
//e.g. foreign keys, duplicates, etc
func AddVariableMappings(context echo.Context) error {

	log.Printf("[handlers] binding new variable mapping...")

	var mappings ac.Batch
	err := context.Bind(&mappings)
	if err != nil {
		msg := fmt.Sprintf("unable to bind JSON to struct: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	lastInserted, err := ac.AddMappings(VAR, VMAP, VID, &mappings)
	if err != nil {
		msg := fmt.Sprintf("variables not added: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	entries, err := ac.GetVariableMappings(lastInserted)
	if err != nil {
		msg := fmt.Sprintf("new entries not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusInternalServerError, msg)
	}

	return context.JSON(http.StatusOK, entries)
}

func AddVariableDefinition(context echo.Context) error {

	log.Printf("[handlers] binding new variable definition...")

	var variable ac.Definition
	err := context.Bind(&variable)
	if err != nil {
		msg := fmt.Sprintf("unable to bind JSON to struct: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	log.Printf("[handlers] adding variable definition...")

	err = ac.AddDefinition(VAR, &variable)
	if err != nil {
		msg := fmt.Sprintf("variable definition failed: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, variable)
}

func EditVariableDefinition(context echo.Context) error {

	log.Printf("[handlers] binding variable definition...")

	var variable ac.Definition
	err := context.Bind(&variable)
	if err != nil {
		msg := fmt.Sprintf("unable to bind JSON to struct: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	log.Printf("[handlers] editing variable definition...")

	err = ac.EditDefinition(VAR, &variable)
	if err != nil {
		msg := fmt.Sprintf("edit failed: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, variable)
}
