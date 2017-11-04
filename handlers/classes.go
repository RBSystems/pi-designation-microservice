package handlers

import (
	"fmt"
	"log"
	"net/http"

	ac "github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

const CLASS_TABLE_NAME = "class_definitions"

func AddClassDefinition(context echo.Context) error {

	log.Printf("[handlers] binding new class definition...")

	var class ac.Definition
	err := context.Bind(&class)
	if err != nil {
		msg := fmt.Sprintf("unable to bind JSON to struct: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	err = ac.AddDefinition(CLASS_TABLE_NAME, &class)
	if err != nil {
		msg := fmt.Sprintf("accessor error: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, class)
}

func EditClassDefinition(context echo.Context) error {

	log.Printf("[handlers] editing class definition...")

	var class ac.Definition
	err := context.Bind(&class)
	if err != nil {
		msg := fmt.Sprintf("unable to bind JSON to struct: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	err = ac.EditDefinition(CLASS_TABLE_NAME, &class)
	if err != nil {
		msg := fmt.Sprintf("accessor error: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, class)
}

func GetClassDefinitionById(context echo.Context) error {

	id, err := ExtractId(context)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	log.Printf("[handlers] fetching class with id: %s", id)

	var class ac.Definition
	err = ac.GetDefinitionById(CLASS_TABLE_NAME, id, &class)
	if err != nil {
		msg := fmt.Sprintf("accessor error: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, class)
}

func GetAllClassDefinitions(context echo.Context) error {

	log.Printf("[handlers] fetching all class definitions")

	var classes []ac.Definition
	err := ac.GetAllDefinitions(CLASS_TABLE_NAME, &classes)
	if err != nil {
		msg := fmt.Sprintf("accessor error: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, classes)
}
