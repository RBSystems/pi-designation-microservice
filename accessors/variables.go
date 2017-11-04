package accessors

import (
	"errors"
	"fmt"
	"log"

	db "github.com/byuoitav/pi-designation-microservice/database"
	"github.com/fatih/color"
)

func GetVariableMappingsById(IDs []int64) ([]VariableMapping, error) {

	log.Printf("[accessors] getting microservice entries...")

	var output []VariableMapping
	for _, id := range IDs {

		var mapping VariableMapping
		err := GetVariableMappingById(id, &mapping)
		if err != nil {
			msg := fmt.Sprintf("entry not found: %s", err.Error())
			log.Printf("%s", color.HiRedString("[accessors] %s", msg))
			return []VariableMapping{}, errors.New(msg)
		}

		output = append(output, mapping)
	}

	return output, nil
}

func GetAllVariableMappings() ([]VariableMapping, error) {

	log.Printf("[accessors] getting all variable mappings...")

	var mappings []DBVariable
	err := db.DB().Select(&mappings, "SELECT * FROM variable_mappings")
	if err != nil {
		msg := fmt.Sprintf("mappings not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return []VariableMapping{}, errors.New(msg)
	}

	var output []VariableMapping

	for _, mapping := range mappings {

		var variable VariableMapping
		err = FillVariableMapping(&mapping, &variable)
		if err != nil {
			msg := fmt.Sprintf("variable not found: %s", err.Error())
			log.Printf("%s", color.HiRedString("[accessors] %s", msg))
			return []VariableMapping{}, errors.New(msg)
		}

		output = append(output, variable)
	}

	return output, nil
}

func GetVariableMappingById(entryID int64, variable *VariableMapping) error {

	log.Printf("[accessors] getting variable entry with ID %d...", entryID)

	//get the IDs
	var mapping DBVariable
	err := db.DB().Get(&mapping, "SELECT * FROM variable_mappings WHERE id = ?", entryID)
	if err != nil {
		msg := fmt.Sprintf("failed to execute query: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	return nil
}

func FillVariableMapping(entry *DBVariable, mapping *VariableMapping) error {

	class, desig, err := GetClassAndDesignation(entry.ClassID, entry.DesigID)
	if err != nil {
		msg := fmt.Sprintf("entry not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	var variable Variable
	err = db.DB().Get(&variable, "SELECT * from variable_definitions WHERE id = ?", entry.VarID)
	if err != nil {
		msg := fmt.Sprintf("entry not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	placeHolder := Mapping{
		ID:          mapping.ID,
		Class:       class,
		Designation: desig,
	}

	mapping.Mapping = placeHolder
	mapping.Variable = variable
	mapping.Value = entry.Value

	return nil
}
