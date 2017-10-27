package handlers

import (
	"fmt"
	"log"
	"net/http"

	ac "github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

const desig = "designation_definitions"

func AddDesignationDefinition(context echo.Context) error {

	log.Printf("[handlers] adding new desigation definition")

	var designation ac.Definition
	err := context.Bind(&designation)
	if err != nil {
		msg := fmt.Sprintf("unable to bind JSON to struct: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusInternalServerError, msg)
	}

	err = ac.AddDefinition(desig, &designation)
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

	err = ac.EditDefinition(desig, &designation)
	if err != nil {
		msg := fmt.Sprintf("entry not updated: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, designation)
}
