package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

func AddRoom(context echo.Context) error {

	log.Printf("[handlers] adding new room...")

	var room accessors.Room
	err := context.Bind(&room)
	if err != nil {
		msg := fmt.Sprintf("unable to bind JSON to struct: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", err.Error()))
		return context.JSON(http.StatusBadRequest, msg)
	}

	err = accessors.AddRoom(&room)
	if err != nil {
		log.Printf("%s", color.HiRedString("[handlers] %s", err.Error()))
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	log.Printf("%s", color.HiGreenString("[accessors] successfully added new room: %s", room.Name))
	return context.JSON(http.StatusOK, room)
}

func EditRoom(context echo.Context) error {

	log.Printf("[handlers] editing room...")

	var room accessors.Room
	err := context.Bind(&room)
	if err != nil {
		msg := fmt.Sprintf("unable to bind JSON to struct: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", err.Error()))
		return context.JSON(http.StatusBadRequest, msg)
	}

	err = accessors.EditRoom(&room)
	if err != nil {
		log.Printf("%s", color.HiRedString("[handlers] %s", err.Error()))
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	log.Printf("%s", color.HiGreenString("[accessors] successfully applied edits to room: %s", room.Name))
	return context.JSON(http.StatusOK, room)

}

func GetAllRooms(context echo.Context) error {

	log.Printf("[handlers] fetching all rooms...")

	rooms, err := accessors.GetAllRooms()
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, rooms)
}

func GetRoomById(context echo.Context) error {

	id := context.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		return context.JSON(http.StatusBadRequest, fmt.Sprintf("invalid ID: %s", err.Error()))
	}

	log.Printf("[handlers] fetching room by ID: %d...", intId)

	room, err := accessors.GetRoomById(int64(intId))
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	return context.JSON(http.StatusOK, room)
}

func DeleteRoom(context echo.Context) error {
	id := context.Param("id")
	intId, err := strconv.Atoi(id)
	if err != nil {
		return context.JSON(http.StatusBadRequest, fmt.Sprintf("invalid ID: %s", err.Error()))
	}

	log.Printf("[handlers] deleting room with ID: %s", intId)

	err = accessors.DeleteRoom(int64(intId))
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	return context.JSON(http.StatusOK, "room successfully deleted")

}
