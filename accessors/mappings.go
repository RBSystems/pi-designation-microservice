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
func AddMappings(mapTable, colName, defId, string, entries *Batch) error {

	if len(entries.Value) == 0 {
		msg := "invalid variable value"
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	//we're building one giant commit
	tx, err := db.DB().Beginx()
	if err != nil {
		msg := fmt.Sprintf("could not begin transaction: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	for class, designations := range entries.Classes {

		for _, designation := range designations {

			//format SQL
			command := fmt.Sprintf("INSERT INTO %s (%s, designation_id, class_id, %s) VALUES (?, ?, ?, ?)", mapTable, defId, colName)

			log.Printf("[accessors] SQL: %s", command)
			_, err = tx.Exec(command, entries.ID, designation, class, entries.Value)
		}
	}

	err = tx.Commit()
	if err != nil {
		msg := fmt.Sprintf("unable to prepare statement: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	return nil
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

func GetVariableMapping(entryID int64) (VariableMapping, error) {

	log.Printf("[accessors] getting microservice entry...")

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
