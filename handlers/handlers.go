package handlers

import (
	"log"
	"net/http"

	"github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

func GetDevice(context echo.Context) error {

	log.Printf("[handlers] getting device...")

	return context.JSON(http.StatusOK, "")
}

func GetEnvironmentVariables(context echo.Context) error {

	host := context.Param("host")
	log.Printf("[handlers] getting environment variables from %s...", host)

	device, err := accessors.GetEnv(host)
	if err != nil {
		log.Printf("%s", color.HiRedString("[error] %s", err.Error))
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	return context.JSON(http.StatusOK, device)
}

func GetUiConfig(context echo.Context) error {

	host := context.Param("host")
	log.Printf("[handlers] getting ui configuration of %s...", host)

	config, err := accessors.GetUi(host)
	if err != nil {
		log.Printf("%s", color.HiRedString("[error] %s", err.Error()))
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	return context.JSON(http.StatusOK, config)
}

func AddNewDevice(context echo.Context) error {

	//host := context.Param("host")
	//	designation := context.Param("designation")

	log.Printf("[handlers] adding new %s device %s...")

	return context.JSON(http.StatusOK, "")
}

func AddNewRoom(context echo.Context) error {

	return nil
}
