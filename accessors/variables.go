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

	raw := "INSERT INTO variable_mappings (variable_id, designation_id, class_id, value) SELECT variable_definitions.id as 'variable', designation_definitions.id as 'designation', class_definitions.id as 'class', :value as 'value' FROM variable_definitions  JOIN designation_definitions on designation_definitions.name=:designation JOIN class_definitions on class_definitions.name=:class  WHERE variable_definitions.name=:name"

	//address each class
	for class, designations := range mappings.Classes {

		for _, designation := range designations {

			result, err := tx.NamedExec(raw, map[string]interface{}{
				"name":        mappings.Name,
				"value":       mappings.Value,
				"class":       class,
				"designation": designation,
			})
			if err != nil {
				msg := fmt.Sprintf("invalid named statement: %s", err.Error())
				log.Printf("%s", color.HiRedString("[accessors] %s", msg))
				tx.Rollback()
				return []Variable{}, errors.New(msg)
			}

			id, err := result.LastInsertId()
			if err != nil { //if it only errors out here, well...
				msg := fmt.Sprintf("last inserted ID not found: %s", err.Error())
				log.Printf("%s", color.HiRedString("[accessors] %s", msg))
				return []Variable{}, errors.New(msg)
			}

			lastInserted = append(lastInserted, Variable{
				ID:    id,
				Value: mappings.Value,
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
