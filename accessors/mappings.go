package accessors

import (
	"errors"
	"fmt"
	"log"

	db "github.com/byuoitav/pi-designation-microservice/database"
	"github.com/fatih/color"
)

//we're assuming the user knows the IDs for everything
//mapTable - name of table to add entries to
//colName - name of column in table to add entries to
//defId - name of column in table to add external ID to
//returns a slice of newly created IDs
func AddMappings(mapTable, colName, defId string, entries *Batch) ([]int64, error) {

	if len(entries.Value) == 0 {
		msg := "invalid mapping value"
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return []int64{}, errors.New(msg)
	}

	//we're building one giant commit
	tx, err := db.DB().Beginx()
	if err != nil {
		msg := fmt.Sprintf("failed to begin transaction: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return []int64{}, errors.New(msg)
	}

	var output []int64

	for _, designations := range entries.Classes {

		for _, designation := range designations.Designations {

			//format SQL
			command := fmt.Sprintf("INSERT INTO %s (%s, designation_id, class_id, %s) VALUES (?, ?, ?, ?)", mapTable, defId, colName)

			log.Printf("[accessors] SQL: %s", command)
			result, err := tx.Exec(command, entries.ID, designation, designations.ID, entries.Value)
			if err != nil {
				msg := fmt.Sprintf("insert action failed: %s", err.Error())
				log.Printf("%s", color.HiRedString("[accessors] %s", msg))
				log.Printf("%s", color.HiRedString("[accessors] rolling back..."))
				tx.Rollback()
				return []int64{}, errors.New(msg)
			}

			id, err := result.LastInsertId()
			if err != nil {
				msg := fmt.Sprintf("last inserted ID not found: %s", err.Error())
				log.Printf("%s", color.HiRedString("[accessors] %s", msg))
				log.Printf("%s", color.HiRedString("[accessors] rolling back..."))
				tx.Rollback()
				return []int64{}, errors.New(msg)
			}

			output = append(output, id)

		}
	}

	err = tx.Commit()
	if err != nil {
		msg := fmt.Sprintf("commit failed: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return []int64{}, errors.New(msg)
	}

	return output, nil
}

func GetMicroserviceMappings(IDs []int64) ([]MicroserviceMapping, error) {

	log.Printf("[accessors] getting microservice entries...")

	var output []MicroserviceMapping
	for _, id := range IDs {

		mapping, err := GetMicroserviceMapping(id)
		if err != nil {
			msg := fmt.Sprintf("entry not found: %s", err.Error())
			log.Printf("%s", color.HiRedString("[accessors] %s", msg))
			return []MicroserviceMapping{}, errors.New(msg)
		}

		output = append(output, mapping)
	}

	return output, nil
}

func GetMicroserviceMapping(entryID int64) (MicroserviceMapping, error) {

	log.Printf("[accessors] getting microservice entry...")

	//get the IDs
	var mapping DBMicroservice
	err := db.DB().Get(&mapping, "SELECT * FROM microservice_mappings WHERE id = ?", entryID)
	if err != nil {
		msg := fmt.Sprintf("failed to execute query: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return MicroserviceMapping{}, errors.New(msg)
	}

	//TODO:make sure it's not the empty set
	//does Get() take care of that?

	class, desig, err := GetClassAndDesignation(mapping.ClassID, mapping.DesigID)
	if err != nil {
		msg := fmt.Sprintf("entry not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return MicroserviceMapping{}, errors.New(msg)
	}

	var microservice Microservice
	err = db.DB().Get(&microservice, "SELECT * from microservice_definitions WHERE id = ?", mapping.MicroID)
	if err != nil {
		msg := fmt.Sprintf("entry not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return MicroserviceMapping{}, errors.New(msg)
	}

	placeHolder := Mapping{
		ID:          mapping.ID,
		Class:       class,
		Designation: desig,
	}

	return MicroserviceMapping{
		Mapping:      placeHolder,
		Microservice: microservice,
		YAML:         mapping.YAML,
	}, nil

}

func GetVariableMappings(IDs []int64) ([]VariableMapping, error) {

	log.Printf("[accessors] getting microservice entries...")

	var output []VariableMapping
	for _, id := range IDs {

		mapping, err := GetVariableMapping(id)
		if err != nil {
			msg := fmt.Sprintf("entry not found: %s", err.Error())
			log.Printf("%s", color.HiRedString("[accessors] %s", msg))
			return []VariableMapping{}, errors.New(msg)
		}

		output = append(output, mapping)
	}

	return output, nil
}

func GetVariableMapping(entryID int64) (VariableMapping, error) {

	log.Printf("[accessors] getting variable entry...")

	//get the IDs
	var mapping DBVariable
	err := db.DB().Get(&mapping, "SELECT * FROM microservice_mappings WHERE id = ?", entryID)
	if err != nil {
		msg := fmt.Sprintf("failed to execute query: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return VariableMapping{}, errors.New(msg)
	}

	class, desig, err := GetClassAndDesignation(mapping.ClassID, mapping.DesigID)
	if err != nil {
		msg := fmt.Sprintf("entry not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return VariableMapping{}, errors.New(msg)
	}

	var variable Variable
	err = db.DB().Get(&variable, "SELECT * from variable_definitions WHERE id = ?", mapping.VarID)
	if err != nil {
		msg := fmt.Sprintf("entry not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return VariableMapping{}, errors.New(msg)
	}

	placeHolder := Mapping{
		ID:          mapping.ID,
		Class:       class,
		Designation: desig,
	}
	return VariableMapping{
		Mapping:  placeHolder,
		Variable: variable,
		Value:    mapping.Value,
	}, nil

}

func GetClassAndDesignation(classID, designationID int64) (class Class, designation Designation, err error) {

	err = db.DB().Get(&class, "SELECT * from class_definitions WHERE id = ?", classID)
	if err != nil {
		err = errors.New(fmt.Sprintf("class not found: %s", err.Error()))
		return
	}

	err = db.DB().Get(&designation, "SELECT * from designation_definitions WHERE id = ?", designationID)
	if err != nil {
		err = errors.New(fmt.Sprintf("designation not found: %s", err.Error()))
		return
	}

	return
}
