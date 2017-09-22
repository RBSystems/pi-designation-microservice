package handlers

import (
	"fmt"
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

	name := context.Param("room")
	designation := context.Param("designation")
	log.Printf("[handlers] adding new %s room %s", designation, name)

	var room accessors.Room

	err := context.Bind(&room)
	if err != nil {
		msg := fmt.Sprintf("unable to unmarshal JSON object: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	room, err = accessors.AddNewRoom(room, designation)
	if err != nil {
		msg := fmt.Sprintf("unable to add new room: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusInternalServerError, msg)
	}

	return context.JSON(http.StatusOK, room)
}
