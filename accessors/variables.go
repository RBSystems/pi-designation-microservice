package accessors

import (
	"database/sql"
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
	ID    int         `json:"id",omitempty`
}

func ValidateVar(variable Variable) error {

	log.Printf("[accessors] validating variable: %s", variable.Key)

	if (len(variable.Key) == 0) || (len(variable.Value) == 0) {
		msg := "empty key or value"
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	return nil
}

func AddNewVariable(variable Variable) error {

	log.Printf("[accessors] adding variable %s: %s with designation %s...", variable.Key, variable.Value, variable.Desig.Name)

	_, err := database.DB().Exec(`INSERT into variables (designation_ID, variable_key, variable_value) values(?,?,?)`, variable.Desig.ID, variable.Key, variable.Value)
	if err != nil {
		msg := fmt.Sprintf("unable to add row to table: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	return nil
}

//given a variable key and a designation ID, this fills the value and the ID of a variable
//@pre: variable has a valid designation field and name
func FillVariable(variable *Variable) error {

	log.Printf("[accessors] searching for %s %s", variable.Desig.Name, variable.Key)
	log.Printf("%s", color.HiBlueString("variable: %v", variable))
	log.Printf("%s", color.HiGreenString("%s", variable.Key))
	log.Printf("%s", color.HiGreenString("%d", variable.Desig.ID))
	log.Printf("%s", color.HiGreenString("%d", variable.ID))
	log.Printf("%s", color.HiGreenString("%s", variable.Value))

	rows, err := database.DB().Query("SELECT variable_ID, variable_value from variables WHERE variable_key = ? AND designation_ID = ?", variable.Key, variable.Desig.ID)
	if err != nil {
		msg := fmt.Sprintf("unable to get row data for variable: %s", err.Error())
		return errors.New(msg)
	}

	defer rows.Close()

	for rows.Next() {
		rows.Scan(&variable.ID, &variable.Value)
	}

	return nil
}

//given a varialbe key and a designation ID, this updates the database to reflect what's in the struct
//@pre variable has a valid designation
func EditVariable(variable Variable) error {

	log.Printf("[accessors] updating %s %s...", variable.Desig.Name, variable.Key)

	result, err := database.DB().Exec("UPDATE variables SET variable_value = ? WHERE variable_key = ? AND designation_ID = ?", variable.Value, variable.Key, variable.Desig.ID)
	if err != nil {
		msg := fmt.Sprintf("unable to update row: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		msg := fmt.Sprintf("unknown number of rows affected: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	log.Printf("[accessors] rows affected: %d", rows)

	if rows == 0 {
		msg := fmt.Sprintf("no rows found with key: %s and designation ID: %d", variable.Key, variable.Desig.ID)
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	return nil
}

//given a variable key and a designation ID, this deletes the variable from the database
//@pre variable has a vaild designation
func DeleteVariable(variable Variable) error {

	log.Printf("[accessors] removing %s %s from database...", variable.Desig.Name, variable.Key)

	_, err := database.DB().Query("DELETE from variables WHERE variable_key = ? AND designation_ID = ?", variable.Key, variable.Desig.ID)
	if err != nil {
		msg := fmt.Sprintf("unable to delete row: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors %s", msg))
		return errors.New(msg)
	}

	return nil
}

//returns a dump of all variables
//smelly
func GetAllVariables() ([]Variable, error) {

	log.Printf("[accessors] fetching all variables...")

	rows, err := database.DB().Query("SELECT * from variables")
	if err != nil {
		msg := fmt.Sprintf("unable to get rows: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors %s", msg))
		return []Variable{}, errors.New(msg)
	}

	return extractVariableRows(rows)
}

//returns a dump of all the variables with the given designation
func GetVariablesByDesignation(designation Designation) ([]Variable, error) {

	log.Printf("[accessors] fetching all variables with designation: %s", designation.Name)

	rows, err := database.DB().Query("SELECT * from variables WHERE designation_ID = ?", designation.ID)
	if err != nil {
		msg := fmt.Sprintf("unable to get rows: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors %s", msg))
		return []Variable{}, errors.New(msg)
	}

	return extractVariableRows(rows)

}

func extractVariableRows(rows *sql.Rows) ([]Variable, error) {

	log.Printf("extracting row data...")

	defer rows.Close()

	var output []Variable
	var variable Variable
	for rows.Next() {

		err := rows.Scan(&variable.ID, &variable.Desig.ID, &variable.Key, &variable.Value)
		if err != nil {
			msg := fmt.Sprintf("unable scan row: %s", err.Error())
			log.Printf("%s", color.HiRedString("[accessors %s", msg))
			return []Variable{}, errors.New(msg)
		}

		variable.Desig, err = GetDesignationById(variable.Desig.ID)
		if err != nil {
			msg := fmt.Sprintf("unable get designation: %s", err.Error())
			log.Printf("%s", color.HiRedString("[accessors %s", msg))
			return []Variable{}, errors.New(msg)
		}

		output = append(output, variable)

	}

	return output, nil
}
