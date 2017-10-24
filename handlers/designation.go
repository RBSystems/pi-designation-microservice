package handlers

import (
	"fmt"
	"log"
	"net/http"

	ac "github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

func AddDesignationDefinition(context echo.Context) error {

	log.Printf("[handlers] adding new desigation definition")

	var designation ac.Designation
	err := context.Bind(&designation)
	if err != nil {
		msg := fmt.Sprintf("unable to bind JSON to struct %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusInternalServerError, msg)
	}

	err = ac.AddDesignationDefinition(&designation)
	if err != nil {
		msg := fmt.Sprintf("error adding designation: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusInternalServerError, msg)
	}

	log.Printf("%s", color.HiGreenString("[handlers] successfully added desigation: %s", designation.Name))

	return context.JSON(http.StatusOK, designation)
}

//this will cause a cascading delete
//be careful doing this
func DeleteDesignation(context echo.Context) error {

	desig := context.Param("designation")
	log.Printf("[handlers] removing designation definition")

	designation := ac.Designation{Name: desig}
	err := ac.DeleteDesignation(designation)
	if err != nil {
		msg := fmt.Sprintf("designation not removed: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusInternalServerError, msg)
	}

	return context.JSON(http.StatusOK, designation)
}
