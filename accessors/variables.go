package accessors

import (
	"errors"
	"fmt"
	"log"

	db "github.com/byuoitav/pi-designation-microservice/database"
	"github.com/fatih/color"
)

func AddVariableDefinition(variable *VariableDefinition) error {

	if len(variable.Name) == 0 {
		msg := "invalid variable name"
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	log.Printf("[accessors] adding new variable definition: %s", variable.Name)

	varDef := "INSERT INTO variable_definitions (name, description) VALUES (?, ?)"
	result, err := db.DB().Exec(varDef, variable.Name, variable.Description)
	if err != nil {
		log.Printf("%s", color.HiRedString("[accessors] %s", err.Error()))
		return err
	}

	variable.ID, err = result.LastInsertId()
	if err != nil {
		log.Printf("%s", color.HiRedString("[accessors] %s", err.Error()))
		return err
	}

	return nil
}

func AddVariableMappings(mappings *VariableBatch) ([]Variable, error) {

	if len(mappings.Value) == 0 {
		msg := "invalid variable value"
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return []Variable{}, errors.New(msg)
	}

	//we're building one giant commit
	tx, err := db.DB().Beginx()
	if err != nil {
		msg := fmt.Sprintf("could not begin transaction: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return []Variable{}, errors.New(msg)
	}

	var lastInserted []Variable
	var desig Designation
	var classDef Class
	var def VariableDefinition

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
				return []Variable{}, errors.New(msg)
			}

			err = db.DB().Get(&desig, "SELECT * FROM designation_definitions WHERE name = ?", designation)
			if err != nil {
				msg := fmt.Sprintf("designation not found: %s", err.Error())
				log.Printf("%s", color.HiRedString("[accessors] %s", msg))
				log.Printf("%s", color.HiRedString("[accessors] rolling back transaction..."))
				tx.Rollback()
				return []Variable{}, errors.New(msg)
			}

			err = db.DB().Get(&def, "SELECT * FROM variable_definitions WHERE name = ?", mappings.Name)
			if err != nil {
				msg := fmt.Sprintf("variable definition not found: %s", err.Error())
				log.Printf("%s", color.HiRedString("[accessors] %s", msg))
				log.Printf("%s", color.HiRedString("[accessors] rolling back transaction..."))
				tx.Rollback()
				return []Variable{}, errors.New(msg)
			}

			result, err := tx.NamedExec("INSERT INTO variable_mappings (value, designation_id, class_id, variable_id) VALUES (:value, :designation, :class, :variable)",
				map[string]interface{}{
					"value":       mappings.Value,
					"designation": desig.ID,
					"class":       classDef.ID,
					"variable":    def.ID,
				})
			if err != nil {
				msg := fmt.Sprintf("failed to add entry: %s", err.Error())
				log.Printf("%s", color.HiRedString("[accessors] %s", msg))
				log.Printf("%s", color.HiRedString("[accessors] rolling back transaction..."))
				tx.Rollback()
				return []Variable{}, errors.New(msg)
			}

			id, err := result.LastInsertId()
			if err != nil {
				msg := fmt.Sprintf("variable ID not found: %s", err.Error())
				log.Printf("%s", color.HiRedString("[accessors] %s", msg))
				log.Printf("%s", color.HiRedString("[accessors] rolling back transaction..."))
				tx.Rollback()
				return []Variable{}, errors.New(msg)
			}

			lastInserted = append(lastInserted, Variable{
				ID:          id,
				Variable:    def,
				Value:       mappings.Value,
				Class:       classDef,
				Designation: desig,
			})

		}
	}

	err = tx.Commit()
	if err != nil {
		msg := fmt.Sprintf("unable to prepare statement: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return []Variable{}, errors.New(msg)
	}

	return lastInserted, nil
}
