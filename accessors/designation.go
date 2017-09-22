package accessors

import (
	"errors"
	"fmt"
	"log"

	"github.com/byuoitav/pi-designation-microservice/database"
	"github.com/fatih/color"
)

func GetDesignationById(ID int) (designation Designation, err error) {

	log.Printf("[accessors] getting room designation by id: %v...", ID)

	err = database.DB().QueryRow(`SELECT * from designation_definition where designation = ?`, ID).Scan(designation.Name, designation.ID)
	if err != nil {
		msg := fmt.Sprintf("unable to query database: %s", err.Error())
		log.Printf("[accessors] %s", color.HiRedString("%s", msg))
		err = errors.New(msg)
		return
	}

	return
}

func GetDesignationByName(name string) (designation Designation, err error) {

	log.Printf("[accessors] getting room designtion by name: %s...", name)

	err = database.DB().QueryRow(`SELECT * from designation_definition where designation_ID = ?`, name).Scan(designation.Name, designation.ID)
	if err != nil {
		msg := fmt.Sprintf("unable to query database: %s", err.Error())
		log.Printf("[accessors] %s", color.HiRedString("%s", msg))
		err = errors.New(msg)
		return
	}

	return
}

func GetAllDesignations() ([]Designation, error) {

	log.Printf("[accessors] getting all room designations...")

	rows, err := database.DB().Query(`SELECT * from designation_definition`)
	if err != nil {
		msg := fmt.Sprintf("unable to execute query %s", err.Error())
		log.Printf("[accessors] %s", color.HiRedString("%s", msg))
		return []Designation{}, errors.New(msg)
	}

	defer rows.Close()

	var designation Designation
	var output []Designation
	for rows.Next() {

		err = rows.Scan(designation.Name, designation.ID)
		if err != nil {
			msg := fmt.Sprintf("unable to scan row: %s", err.Error())
			log.Printf("[accessors] %s", color.HiRedString("%s", msg))
			return []Designation{}, errors.New(msg)
		}

		output = append(output, designation)

	}

	return output, nil

}

func AddDesignation(designation Designation) error {

	log.Printf("[accessors] adding new room designation %s", designation.Name)

	return nil
}

func RemoveDesignation(designation Designation) error {

	log.Printf("[accessors] removing room desigation %s", designation.Name)

	return nil
}

func UpdateDesignation(new, old Designation) error {

	log.Printf("[accessors] updating room designation %s", old.Name)

	return nil
}
