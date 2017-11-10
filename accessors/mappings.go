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
func AddMappings(mappingTable, definitionColumnName, valueColumnName string, entries *Batch) ([]int64, error) {

	if len(entries.Value) == 0 {
		msg := "invalid mapping value"
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return []int64{}, errors.New(msg)
	}

	var output []int64

	for _, class := range entries.Classes {

		for _, designation := range class.Designations {

			id, err := AddMapping(mappingTable, definitionColumnName, valueColumnName, entries.Value, entries.ID, class.ID, designation)
			if err != nil {
				msg := fmt.Sprintf("failed to add single mapping: %s", err.Error())
				log.Printf("%s", color.HiRedString("[accessors] %s", msg))
				return []int64{}, errors.New(msg)
			}

			output = append(output, id)

		}
	}

	return output, nil
}

//@param mappingTable - name of the table to insert mapping into
//@param definitionColumnName - name of the column the definition ID goes into
//@param valueColumnName - name of the column the value is stored in
//@param value - value to be stored
//@param entryID - ID of entry, e.g. variable_id, microservice_id
//@param designationID - designation ID of mapping
//@param classID - ID of class (e.g. av-control)
//the only string value that should come from the user is 'value'
func AddMapping(mappingTable, definitionColumnName, valueColumnName, value string, entryID, classID, designationID int64) (int64, error) {

	log.Printf("[accessors] adding mapping...")

	//format SQL
	command := fmt.Sprintf("INSERT INTO %s (%s, designation_id, class_id, %s) VALUES (?, ?, ?, ?)", mappingTable, definitionColumnName, valueColumnName)
	log.Printf("[accessors] SQL: %s", command)

	result, err := db.DB().Exec(command, entryID, designationID, classID, value)
	if err != nil {
		msg := fmt.Sprintf("insert action failed: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return 0, errors.New(msg)
	}

	id, err := result.LastInsertId()
	if err != nil {
		msg := fmt.Sprintf("last inserted ID not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return 0, errors.New(msg)
	}

	return id, nil
}

func EditMapping(mappingTable, definitionColumnName, valueColumnName, value string, definitionID, classID, designationID, mappingID int64) error {

	log.Printf("[accessors] editing mapping...")

	//format SQL
	command := fmt.Sprintf("UPDATE %s SET %s = ?, class_id = ?, designation_id = ?, %s = ? WHERE id = ?", mappingTable, definitionColumnName, valueColumnName)
	log.Printf("[accessors] SQL: %s", command)

	_, err := db.DB().Exec(command, definitionID, classID, designationID, value, mappingID)
	if err != nil {
		msg := fmt.Sprintf("edit failed: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	return nil
}

func GetAllMicroserviceMappings() ([]MicroserviceMapping, error) {

	log.Printf("[accessors] getting all microservice mappings...")

	var mappings []DBMicroservice
	err := db.DB().Select(&mappings, "SELECT * FROM microservice_mappings")
	if err != nil {
		msg := fmt.Sprintf("mappings not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return []MicroserviceMapping{}, errors.New(msg)
	}

	var output []MicroserviceMapping

	for _, mapping := range mappings {

		var microservice MicroserviceMapping
		err = FillMicroserviceMapping(&mapping, &microservice)
		if err != nil {
			msg := fmt.Sprintf("microservice not found: %s", err.Error())
			log.Printf("%s", color.HiRedString("[accessors] %s", msg))
			return []MicroserviceMapping{}, errors.New(msg)
		}

		output = append(output, microservice)
	}

	return output, nil

}

func GetMicroserviceMappingsById(IDs []int64) ([]MicroserviceMapping, error) {

	log.Printf("[accessors] getting microservice entries...")

	var output []MicroserviceMapping
	for _, id := range IDs {

		var microservice MicroserviceMapping
		err := GetMicroserviceMappingById(id, &microservice)
		if err != nil {
			msg := fmt.Sprintf("entry not found: %s", err.Error())
			log.Printf("%s", color.HiRedString("[accessors] %s", msg))
			return []MicroserviceMapping{}, errors.New(msg)
		}

		output = append(output, microservice)
	}

	return output, nil
}

func GetMicroserviceMappingById(entryID int64, microservice *MicroserviceMapping) error {

	log.Printf("[accessors] getting microservice entry...")

	//get the IDs
	var mapping DBMicroservice
	err := db.DB().Get(&mapping, "SELECT * FROM microservice_mappings WHERE id = ?", entryID)
	if err != nil {
		msg := fmt.Sprintf("failed to execute query: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	//TODO:make sure it's not the empty set
	//does Get() take care of that?

	err = FillMicroserviceMapping(&mapping, microservice)
	if err != nil {
		msg := fmt.Sprintf("failed to execute query: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	return nil

}

//takes entry in the microservice_mappings table and fleshes it out
func FillMicroserviceMapping(mapping *DBMicroservice, output *MicroserviceMapping) error {

	class, desig, err := GetClassAndDesignation(mapping.ClassID, mapping.DesigID)
	if err != nil {
		msg := fmt.Sprintf("entry not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	var microservice Microservice
	err = db.DB().Get(&microservice, "SELECT * from microservice_definitions WHERE id = ?", mapping.MicroID)
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

	output.Mapping = placeHolder
	output.Microservice = microservice
	output.YAML = mapping.YAML

	return nil
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

func DeleteMapping(table string, id int64) error {

	log.Printf("[accessors] deleting entry from table %s with id %d", table, id)

	command := fmt.Sprintf("DELETE FROM %s WHERE id = ?", table)

	_, err := db.DB().Exec(command, id)
	if err != nil {
		msg := fmt.Sprintf("unable to delete mapping: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	return nil
}

func GetDockerComposeByDesignationAndClass(microservices *[]DBMicroservice, classId, desigId int64) error {

	log.Printf("[accessors] querying database for microservice mappings with class ID %d and designation ID %d", classId, desigId)

	err := db.DB().Select(microservices, "SELECT * FROM microservice_mappings WHERE designation_id = ? AND class_id = ?", desigId, classId)
	if err != nil {
		return err
	}

	return nil

}
