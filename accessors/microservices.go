package accessors

import (
	"errors"
	"fmt"
	"log"

	db "github.com/byuoitav/pi-designation-microservice/database"
	"github.com/fatih/color"
)

func AddMicroserviceDefinition(microservice *MicroserviceDefinition) error {

	if len(microservice.Name) == 0 {
		msg := "invalid microservice name"
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	log.Printf("[accessors] adding new microservice: %s", microservice.Name)

	microDef := "INSERT INTO microservice_definitions (name, description) VALUES (?, ?)"
	result, err := db.DB().Exec(microDef, microservice.Name, microservice.Description)
	if err != nil {
		log.Printf("%s", color.HiRedString("[accessors] %s", err.Error()))
		return err
	}

	microservice.ID, err = result.LastInsertId()
	if err != nil {
		log.Printf("%s", color.HiRedString("[accessors] %s", err.Error()))
		return err
	}

	return nil
}

func AddMicroserviceMappings(mappings *MicroserviceBatch) ([]Microservice, error) {

	log.Printf("[accessors] adding new microservice mappings...")

	if len(mappings.YAML) == 0 {
		msg := "invalid variable value"
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return []Microservice{}, errors.New(msg)
	}

	//we're building one giant commit
	tx, err := db.DB().Beginx()
	if err != nil {
		msg := fmt.Sprintf("could not begin transaction: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return []Microservice{}, errors.New(msg)
	}

	var lastInserted []Microservice
	var desig Designation
	var classDef Class
	var def MicroserviceDefinition

	//address each class
	for class, designations := range mappings.Classes {

		for _, designation := range designations {

			//get class definiton
			err = db.DB().Get(&classDef, "SELECT * FROM class_definitions WHERE name = ?", class)
			if err != nil {
				msg := fmt.Sprintf("class not found: %s", err.Error())
				log.Printf("%s", color.HiRedString("[accessors] %s", msg))
				log.Printf("%s", color.HiRedString("[accessors] rolling back transaction..."))
				tx.Rollback()
				return []Microservice{}, errors.New(msg)
			}

			err = db.DB().Get(&desig, "SELECT * FROM designation_definitions WHERE name = ?", designation)
			if err != nil {
				msg := fmt.Sprintf("designation not found: %s", err.Error())
				log.Printf("%s", color.HiRedString("[accessors] %s", msg))
				log.Printf("%s", color.HiRedString("[accessors] rolling back transaction..."))
				tx.Rollback()
				return []Microservice{}, errors.New(msg)
			}

			err = db.DB().Get(&def, "SELECT * FROM microservice_definitions WHERE name = ?", mappings.Name)
			if err != nil {
				msg := fmt.Sprintf("variable definition not found: %s", err.Error())
				log.Printf("%s", color.HiRedString("[accessors] %s", msg))
				log.Printf("%s", color.HiRedString("[accessors] rolling back transaction..."))
				tx.Rollback()
				return []Microservice{}, errors.New(msg)
			}

			result, err := tx.NamedExec("INSERT INTO microservice_mappings (value, designation_id, class_id, microservice_id) VALUES (:value, :designation, :class, :yaml)",
				map[string]interface{}{
					"value":       mappings.YAML,
					"designation": desig.ID,
					"class":       classDef.ID,
					"yaml":        def.ID,
				})
			if err != nil {
				msg := fmt.Sprintf("failed to add entry: %s", err.Error())
				log.Printf("%s", color.HiRedString("[accessors] %s", msg))
				log.Printf("%s", color.HiRedString("[accessors] rolling back transaction..."))
				tx.Rollback()
				return []Microservice{}, errors.New(msg)
			}

			id, err := result.LastInsertId()
			if err != nil {
				msg := fmt.Sprintf("variable ID not found: %s", err.Error())
				log.Printf("%s", color.HiRedString("[accessors] %s", msg))
				log.Printf("%s", color.HiRedString("[accessors] rolling back transaction..."))
				tx.Rollback()
				return []Microservice{}, errors.New(msg)
			}

			lastInserted = append(lastInserted, Microservice{
				ID:           id,
				Microservice: def,
				YAML:         mappings.YAML,
				Class:        classDef,
				Designation:  desig,
			})

		}
	}

	err = tx.Commit()
	if err != nil {
		msg := fmt.Sprintf("unable to prepare statement: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return []Microservice{}, errors.New(msg)
	}

	return lastInserted, nil
}
