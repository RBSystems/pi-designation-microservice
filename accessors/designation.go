package accessors

import (
	"errors"
	"fmt"
	"log"

	db "github.com/byuoitav/pi-designation-microservice/database"
	"github.com/fatih/color"
)

func GetDesignationByName(name string) (Designation, error) {

	log.Printf("[accessors] getting room designation by name: %s...", name)

	var output Designation
	err := db.DB().QueryRow(`SELECT * from designation_definition where designation = ?`, name).Scan(&output.Name, &output.ID)
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

	rows, err := db.DB().Query(`SELECT * from designation_definition`)
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

func AddDesignationDefinition(designation *Designation) error {

	if len(designation.Name) == 0 {
		msg := "invalid designation name"
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	log.Printf("[accessors] adding new designation: %s", designation.Name)

	parentInsert := "INSERT INTO designation_definitions (name, description) VALUES (?, ?)"
	result, err := db.DB().Exec(parentInsert, designation.Name, designation.Description)
	if err != nil {
		msg := fmt.Sprintf("designation not added: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	designation.ID, err = result.LastInsertId()
	if err != nil {
		msg := fmt.Sprintf("id not found: %s", err.Error())
		log.Printf("[accessors] %s", color.HiRedString("%s", msg))
		return errors.New(msg)
	}

	return nil
}

func DeleteDesignation(designation *Designation) error {

	log.Printf("[accessors] removing room desigation %s", designation.Name)

	_, err := db.DB().Exec("DELETE from designation_definition WHERE designation = ?", designation.Name)
	if err != nil {
		msg := fmt.Sprintf("problem deleting designation: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	return nil
}

func EditDesignationDefinition(old, new *Designation) error {

	log.Printf("[handlers] updating designation definition...")

	if len(new.Name) == 0 {
		msg := "invalid designation name"
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return errors.New(msg)
	}

	if len(new.Description) == 0 {
		msg := "invalid description"
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return errors.New(msg)
	}

	result, err := db.DB().Exec("UPDATE designation_definitions SET name=?, description=? WHERE name=?", new.Name, new.Description, old.Name)
	if err != nil {
		msg := fmt.Sprintf("unable to update designation: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	numRows, err := result.RowsAffected()
	if err != nil {
		msg := fmt.Sprintf("ID not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	if numRows < 1 {
		msg := "designation definition to edit not found."
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	id, err := result.LastInsertId()
	if err != nil {
		msg := fmt.Sprintf("ID not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	new.ID = id

	return nil
}
