package handlers

import (
	"log"
	"net/http"

	"github.com/byuoitav/pi-designation-microservice/dbo"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

func GetEnvironmentVariables(context echo.Context) error {

	host := context.Param("host")
	log.Printf("[handlers] getting environment variables from %s...", host)

	device, err := dbo.GetEnv(host)
	if err != nil {
		log.Printf("%s", color.HiRedString("[error] %s", err.Error))
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	return context.JSON(http.StatusOK, device)
}

func GetUiConfig(context echo.Context) error {

	host := context.Param("host")
	log.Printf("[handlers] getting ui configuration of %s...", host)

	config, err := dbo.GetUi(host)
	if err != nil {
		log.Printf("%s", color.HiRedString("[error] %s", err.Error()))
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	return context.JSON(http.StatusOK, config)
}
