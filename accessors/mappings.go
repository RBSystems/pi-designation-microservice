package accessors

import (
	"errors"
	"fmt"
	"log"

	db "github.com/byuoitav/pi-designation-microservice/database"
	"github.com/fatih/color"
)

type Batch struct {
	Name    string              `json:"name"`
	Classes map[string][]string `json:"classes"` //maps a class to a list of designations
	Value   string              `json:"value"`
}

type Mapping struct {
	ID          int64      `json:"id"`
	Type        Definition `json:"type"`        //either microservice or variable definition ...today
	Class       Definition `json:"class"`       //classes, e.g. av-control, scheduling, etc.
	Designation Definition `json:"designation"` //designations exist inside classes
	Value       string     `json:"value"`       //the actual value, e.g. some YAML or an environement variable
}

func AddMappings(defTable, mapTable string, entries *Batch) ([]Mapping, error) {

	if len(entries.Value) == 0 {
		msg := "invalid variable value"
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return []Mapping{}, errors.New(msg)
	}

	//we're building one giant commit
	tx, err := db.DB().Beginx()
	if err != nil {
		msg := fmt.Sprintf("could not begin transaction: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return []Mapping{}, errors.New(msg)
	}

	var lastInserted []Mapping
	var desig, classDef, typeDef Definition

	//address each class
	for class, designations := range entries.Classes {

		for _, designation := range designations {

			//get class definiton
			err = tx.Get(&classDef, "SELECT * FROM class_definitions WHERE name = ?", class)
			if err != nil {
				msg := fmt.Sprintf("class not found: %s", err.Error())
				log.Printf("%s", color.HiRedString("[accessors] %s", msg))
				log.Printf("%s", color.HiRedString("[accessors] rolling back transaction..."))
				tx.Rollback()
				return []Mapping{}, errors.New(msg)
			}

			err = tx.Get(&desig, "SELECT * FROM designation_definitions WHERE name = ?", designation)
			if err != nil {
				msg := fmt.Sprintf("designation not found: %s", err.Error())
				log.Printf("%s", color.HiRedString("[accessors] %s", msg))
				log.Printf("%s", color.HiRedString("[accessors] rolling back transaction..."))
				tx.Rollback()
				return []Mapping{}, errors.New(msg)
			}

			//format SQL
			command := fmt.Sprintf("SELECT * from %s WHERE name = ?", defTable)
			err = tx.Get(&typeDef, command, entries.Name)
			if err != nil {
				msg := fmt.Sprintf("variable definition not found: %s", err.Error())
				log.Printf("%s", color.HiRedString("[accessors] %s", msg))
				log.Printf("%s", color.HiRedString("[accessors] rolling back transaction..."))
				tx.Rollback()
				return []Mapping{}, errors.New(msg)
			}

			command = fmt.Sprintf("INSERT INTO %s (value, designation_id, class_id, variable_id) VALUES (:value, :designation, :class, :variable)", mapTable)
			result, err := tx.NamedExec(command,
				map[string]interface{}{
					"value":       entries.Value,
					"designation": desig.ID,
					"class":       classDef.ID,
					"variable":    typeDef.ID,
				})
			if err != nil {
				msg := fmt.Sprintf("failed to add entry: %s", err.Error())
				log.Printf("%s", color.HiRedString("[accessors] %s", msg))
				log.Printf("%s", color.HiRedString("[accessors] rolling back transaction..."))
				tx.Rollback()
				return []Mapping{}, errors.New(msg)
			}

			id, err := result.LastInsertId()
			if err != nil {
				msg := fmt.Sprintf("variable ID not found: %s", err.Error())
				log.Printf("%s", color.HiRedString("[accessors] %s", msg))
				log.Printf("%s", color.HiRedString("[accessors] rolling back transaction..."))
				tx.Rollback()
				return []Mapping{}, errors.New(msg)
			}

			lastInserted = append(lastInserted, Mapping{
				ID:          id,
				Type:        typeDef,
				Value:       entries.Value,
				Class:       classDef,
				Designation: desig,
			})

		}
	}

	err = tx.Commit()
	if err != nil {
		msg := fmt.Sprintf("unable to prepare statement: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return []Mapping{}, errors.New(msg)
	}

	return lastInserted, nil
}
