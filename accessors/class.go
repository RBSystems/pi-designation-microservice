package accessors

import (
	"errors"
	"fmt"
	"log"

	"github.com/byuoitav/pi-designation-microservice/database"
	"github.com/fatih/color"
)

func AddClassDefinition(class *Class) error {

	if len(class.Name) == 0 {
		msg := "invalid class name"
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	log.Printf("[accessors] adding new class definition: %s", class.Name)

	classDef := "INSERT INTO class_definitions (name, description) VALUES (?, ?)"
	result, err := database.DB().Exec(classDef, class.Name, class.Description)
	if err != nil {
		msg := fmt.Sprintf("class not added: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	class.ID, err = result.LastInsertId()
	if err != nil {
		msg := fmt.Sprintf("ID not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	return nil
}
