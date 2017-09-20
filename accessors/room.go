package accessors

import (
	"errors"
	"fmt"
	"log"

	"github.com/byuoitav/pi-designation-microservice/database"
	"github.com/fatih/color"
)

func GetEnv(hostname string) (Device, error) {

	log.Printf("[dbo] getting environment variables of hostname: %s", hostname)

	return Device{}, nil
}

func GetUi(hostname string) (string, error) {

	log.Printf("[dbo] getting ui configuration of %s", hostname)

	return "", nil
}

func AddNewRoom(room RoomConfig, designation string) (RoomConfig, error) {

	log.Printf("[dbo] adding new %s room: %s", designation, room.Name)

	log.Printf("Validating designation...")
	err := ValidateDesignation(designation)
	if err != nil {
		msg := fmt.Sprintf("invalid desigation: %s", err.Error())
		log.Printf("[accessors] %s", color.HiRedString("%s", msg))
		return RoomConfig{}, errors.New(msg)
	}

	return room, nil
}

func ValidateDesignation(desig string) error {

	err := database.DB().QueryRow(`SELECT designation_ID from designation_definition where designation = ?`, desig).Scan(
	if err != nil {
		msg := fmt.Sprintf("unable to query database: %s", err.Error())
		log.Printf("[accessors] %s", color.HiRedString("%s", msg))
		return errors.New(msg)
	}

	if len(row) == 0 {
		return errors.New("room designation not found!")
	}

	return nil
}
