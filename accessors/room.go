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

	encoded, err := json.Marshal(room.Config)
	if err != nil {
		msg := fmt.Sprintf("could not marshal room configuration into JSON object", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	//TODO not really sure if we should grab the result? Is it RESTful to get the last inserted ID and the rows affected? Or do we need the entire room struct?
	_, err = database.DB().Exec(`insert into rooms (name, designation_ID, ui_config) values(?, ?, ?)`, room.Name, room.Desig.ID, encoded)
	if err != nil {
		msg := fmt.Sprintf("could not add room to database: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	return nil
}

func GetRoomByName(name string) (Room, error) {

	log.Printf("[accessors] searching for room %s...", name)

	var preConfig string
	var output Room

	err := database.DB().QueryRow(`select rooms.room_ID, rooms.name, designation_definition.designation, rooms.designation_ID, rooms.ui_config from rooms join designation_definition on designation_definition.designation_ID = rooms.designation_ID where rooms.name = ?`, name).Scan(
		&output.ID,
		&output.Name,
		&output.Desig.Name,
		&output.Desig.ID,
		&preConfig)
	if err != nil {
		msg := fmt.Sprintf("could not get room from database %s", err.Error())
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
