package accessors

import (
	"errors"
	"fmt"
	"log"

	db "github.com/byuoitav/pi-designation-microservice/database"
	"github.com/fatih/color"
)

//	variable to be consumed externally
type Variable struct {
	Id             int64  `db:"id"`
	DesignationId  int64  `db:"designation_id"`
	DefinitionId   int64  `db:"variable_id"`
	MicroserviceId int64  `db:"microservice_id"`
	Name           string `db:"name"`
	Value          string `db:"value"`
}

//	defines a variable
//	corresponds to the variable_definitions table
type VariableDefinition struct {
	Id          int64  `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

//	Maps a variable definition to a designation and value
//	Corresponds to variable_mappings table
type VariableMapping struct {
	Id            int64  `db:"id"`
	Value         string `db:"value"`
	DesignationId int64  `db:"designation_id"`
	VariableId    int64  `db:"variable_id"`
}

//	maps microservices and definitions to variables
//	ie microservice 'a' with designation 'b' requires variable 'c'
//	variable 'c' will have the value associated with designation 'b'
//	correpsonds to variable_sets table
type VariableRequirement struct {
	Id             int64 `db:"id"`
	MicroserviceId int64 `db:"microservice_id"`
	DesignationId  int64 `db:"designation_id"`
	VariableId     int64 `db:"variable_id"`
}

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
	return nil
}

func FillVariableMapping(entry *DBVariable, mapping *VariableMapping) error {
	return nil
}

func GetVariablesByClassAndDesignation(classId, desigId int64) ([]VariableMapping, error) {

	log.Printf("[accessors] querying database for variable mappings with class ID %d and designation ID %d", classId, desigId)

	var preMappings []DBVariable
	err := db.DB().Select(&preMappings, "SELECT * FROM variable_mappings WHERE designation_id = ? AND class_id = ?", desigId, classId)
	if err != nil {
		return []VariableMapping{}, err
	}

	var output []VariableMapping
	for _, mapping := range preMappings {

		var variable VariableMapping
		err = FillVariableMapping(&mapping, &variable)
		if err != nil {
			return []VariableMapping{}, err
		}

		output = append(output, variable)
	}

	return output, nil
}

func GetVariablesByMicroserviceAndDesignation(microId, desigId int64) ([]Variable, error) {

	query := `SELECT variable_mappings.id, variable_mappings.variable_id, variable_definitions.name, variable_mappings.designation_id, variable_sets.microservice_id, variable_mappings.value 
				FROM variable_mappings 
				JOIN variable_sets 
					ON variable_sets.designation_id = variable_mappings.designation_id 
					AND variable_mappings.variable_id = variable_sets.variable_id
				JOIN variable_definitions 
					ON variable_definitions.id = variable_mappings.variable_id
				WHERE variable_sets.microservice_id = ? and variable_sets.designation_id = ?`

	var output []Variable

	err := db.DB().Select(&output, query, microId, desigId)
	if err != nil {
		return nil, err
	}

	return output, nil
}
