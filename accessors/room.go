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

func AddNewRoom(room Room) (Room, error) {

	log.Printf("[dbo] adding new %s room: %s", room.Desig.Name, room.Name)

	err := ValidateDesignation(room.Desig)
	if err != nil {
		msg := fmt.Sprintf("invalid desigation: %s", err.Error())
		log.Printf("[accessors] %s", color.HiRedString("%s", msg))
		return Room{}, errors.New(msg)
	}

	return room, nil
}

func ValidateDesignation(desig Designation) error {

	log.Printf("[accessors] validating designation... %s", desig)

	err := database.DB().QueryRow(`SELECT designation_ID from designation_definition where designation = ?`, desig).Scan(desig.Name, desig.ID)
	if err != nil {
		msg := fmt.Sprintf("unable to query database: %s", err.Error())
		log.Printf("[accessors] %s", color.HiRedString("%s", msg))
		return errors.New(msg)
	}

	if (desig.Name == nil) || (desig.ID == nil) {
		return errors.New("invalid designation")
	}

	return nil
}
