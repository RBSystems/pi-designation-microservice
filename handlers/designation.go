package handlers

import (
	"fmt"
	"log"
	"net/http"

	ac "github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

func AddDesignation(context echo.Context) error {

	definition := context.Param("definition")
	log.Printf("[handlers] adding new desigation definition: %s", definition)

	designation := ac.Designation{Name: definition}
	err := ac.AddDesignation(&designation)
	if err != nil {
		msg := fmt.Sprintf("designation not added to database: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusInternalServerError, msg)
	}

	return context.JSON(http.StatusOK, designation)
}
