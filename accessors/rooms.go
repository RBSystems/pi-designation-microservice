package accessors

import (
	"errors"
	"fmt"
	"log"

	"github.com/byuoitav/pi-designation-microservice/database"
	"github.com/byuoitav/touchpanel-ui-microservice/uiconfig"
	"github.com/fatih/color"
)

type Room struct {
	Name          string `json:"name"`
	Id            int64  `json:"id"`
	DesignationId int64  `json:"designation-id"`
	UiConfig      uiconfig.UIConfig
}

func AddRoom(room *Room) error {

	result, err := database.DB().Exec("INSERT INTO rooms (designation_id, ui_configuration, name) VALUES(?, ?, ?)", room.DesignationId, room.UiConfig, room.Name)
	if err != nil {
		msg := fmt.Sprintf("unable to add entry: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("%s", color.HiRedString("[accessors] last inserted ID not found: %s", err.Error()))
	}

	room.Id = id
	return nil

}

func EditRoom(room *Room) error {

	result, err := database.DB().Exec("UPDATE rooms SET designation_id = ?, ui_configuration = ?, name = ? WHERE id = ?")
	if err != nil {
		msg := fmt.Sprintf("unable to edit entry: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	numRows, err := result.RowsAffected()
	if err != nil {
		log.Printf("%s", color.HiRedString("[accessors] number of rows affected not found: %s", err.Error()))
	}

	if numRows == 0 {
		msg := "invalid edit"
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	return nil
}

func GetAllRooms() ([]Room, error) {

	var output []Room
	err := database.DB().Get(&output, "SELECT * FROM rooms")
	if err != nil {
		msg := fmt.Sprintf("entries not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return []Room{}, errors.New(msg)
	}

	return output, nil
}

func GetRoomById(id int64) (Room, error) {

	var output Room
	err := database.DB().Get(&output, "SELECT * FROM rooms WHERE id = ?", id)
	if err != nil {
		msg := fmt.Sprintf("entry not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return Room{}, errors.New(msg)
	}

	return output, nil
}

func DeleteRoom(id int64) error {

	result, err := database.DB().Exec("DELETE FROM rooms WHERE id = ?", id)
	if err != nil {
		msg := fmt.Sprintf("failed to delete entry: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	numRows, err := result.RowsAffected()
	if err != nil {
		log.Printf("%s", color.HiRedString("[accessors] number of rows deleted not found: %s", err.Error()))
	}

	if numRows == 0 {
		msg := "no rows deleted"
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	return nil
}
