package accessors

import (
	"errors"
	"fmt"
	"log"

	"github.com/byuoitav/pi-designation-microservice/database"
	"github.com/fatih/color"
)

func GetDesignationById(ID int64) (designation Designation, err error) {

	log.Printf("[accessors] getting room designation by id: %v...", ID)

	err = database.DB().QueryRow(`SELECT * from designation_definition where designation = ?`, ID).Scan(designation.Name, designation.ID)
	if err != nil {
		msg := fmt.Sprintf("problem with query: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		err = errors.New(msg)
		return
	}

	return
}

func GetDesignationByName(name string) (Designation, error) {

	log.Printf("[accessors] getting room designation by name: %s...", name)

	var output Designation
	err := database.DB().QueryRow(`SELECT * from designation_definition where designation = ?`, name).Scan(&output.Name, &output.ID)
	if err != nil {
		msg := fmt.Sprintf("problem with query: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return Designation{}, errors.New(msg)
	}

	log.Printf("%s", color.HiCyanString("Sheriff, this is no time to panic!"))

	return output, nil
}

func GetAllDesignations() ([]Designation, error) {

	log.Printf("[accessors] getting all room designations...")

	rows, err := database.DB().Query(`SELECT * from designation_definition`)
	if err != nil {
		msg := fmt.Sprintf("unable to execute query %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return []Designation{}, errors.New(msg)
	}

	defer rows.Close()

	var designation Designation
	var output []Designation
	for rows.Next() {

		err = rows.Scan(designation.Name, designation.ID)
		if err != nil {
			msg := fmt.Sprintf("problem with row scan: %s", err.Error())
			log.Printf("[accessors] %s", color.HiRedString("%s", msg))
			return []Designation{}, errors.New(msg)
		}

		output = append(output, designation)

	}

	return output, nil

}

func AddDesignation(designation *Designation) error {

	log.Printf("[accessors] adding new room designation %s", designation.Name)

	result, err := database.DB().Exec("INSERT into designation_definition (designation) values(?)", designation.Name)
	if err != nil {
		msg := fmt.Sprintf("designation not added: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	id, err := result.LastInsertId()
	if err != nil {
		msg := fmt.Sprintf("new designation ID not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	designation.ID = id

	return nil
}

func DeleteDesignation(designation Designation) error {

	log.Printf("[accessors] removing room desigation %s", designation.Name)

	_, err := database.DB().Exec("DELETE from designation_definition WHERE designation = ?", designation.Name)
	if err != nil {
		msg := fmt.Sprintf("problem deleting designation: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	return nil
}

func UpdateDesignation(new, old Designation) error {

	log.Printf("[accessors] updating room designation %s", old.Name)

	return nil
}
