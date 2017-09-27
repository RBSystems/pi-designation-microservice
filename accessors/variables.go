package accessors

import (
	"errors"
	"fmt"
	"log"

	"github.com/byuoitav/pi-designation-microservice/database"
	"github.com/fatih/color"
)

type Variable struct {
	Key   string      `json:"key"`
	Value string      `json:"value"`
	Desig Designation `json:"designation"`
	ID    string      `json:"id",omitempty`
}

func AddNewVariable(variable Variable) error {

	log.Printf("Adding variable %s: %s with designation %s...", variable.Key, variable.Value, variable.Desig.Name)

	_, err := database.DB().Exec(`INSERT into variables designation_ID, variable_key, variable_value, values(?,?,?)`, variable.Desig.ID, variable.Key, variable.Value)
	if err != nil {
		msg := fmt.Sprintf("unable to add row to table: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	return nil
}
