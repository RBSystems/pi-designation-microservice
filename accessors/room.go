package accessors

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

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

func AddNewRoom(room Room) error {

	log.Printf("[dbo] adding new %s room: %s...", room.Desig.Name, room.Name)

	//TODO not really sure if we should grab the result? Is it RESTful to get the last inserted ID and the rows affected? Or do we need the entire room struct?
	_, err := database.DB().Exec(`insert into rooms (name, designation_ID, ui_config) values(?, ?, ?)`, room.Name, room.Desig.ID, room.Config)
	if err != nil {
		msg := fmt.Sprintf("could not add room to database: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	return nil
}

func GetRoomByName(name string) (Room, error) {

	log.Printf("[dbo] searching for room %s...", name)

	var preConfig string
	var output Room

	err := database.DB().QueryRow(`SELECT rooms.room_ID,
										rooms.name,
										designation_definition.designation,
										rooms.designation_ID,
										rooms.ui_config,
										FROM rooms JOIN designation_definition ON rooms.designation_ID = designation_definition.designation_ID
										WHERE rooms.name = '?';`).Scan(
		&output.ID,
		&output.Name,
		&output.Desig.ID,
		&preConfig)
	if err != nil {
		msg := fmt.Sprintf("could not get room from database %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return Room{}, errors.New(msg)
	}

	//decode room config
	decoder := json.NewDecoder(strings.NewReader(preConfig))
	err = decoder.Decode(&output.Config)
	if err != nil {
		msg := fmt.Sprintf("could not unmarshal room configuration: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return Room{}, errors.New(msg)
	}

	return output, nil
}

//func ValidateDesignation(desig Designation) error {
//
//	log.Printf("[accessors] validating designation... %s", desig)
//
//	return nil
//}
//
//func ValidateRoom(room Room) error {
//
//	log.Printf("[accessors] validating room...")
//
//	if (desig.Name == nil) || (desig.ID == nil) {
//		return errors.New("invalid designation")
//	}
//
//}
