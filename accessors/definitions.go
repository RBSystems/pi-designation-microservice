package accessors

import (
	"errors"
	"fmt"
	"log"

	db "github.com/byuoitav/pi-designation-microservice/database"
	"github.com/fatih/color"
)

func AddDefinition(table string, def *Definition) error {

	log.Printf("[accessors] adding definition to %s...", table)

	if len(def.Name) == 0 {
		msg := "invalid definition name"
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	log.Printf("[accessors] adding new definition %s to table %s", def.Name, table)

	insert := fmt.Sprintf("INSERT INTO %s (name, description) VALUES (?, ?)", table)
	result, err := db.DB().Exec(insert, def.Name, def.Description)
	if err != nil {
		msg := fmt.Sprintf("definition not added: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	def.ID, err = result.LastInsertId()
	if err != nil {
		msg := fmt.Sprintf("id not found: %s", err.Error())
		log.Printf("[accessors] %s", color.HiRedString("%s", msg))
		return errors.New(msg)
	}

	return nil
}

func EditDefinition(table string, def *Definition) error {

	log.Printf("[accessors] updating definition in %s...", table)

	//validate input
	if len(def.Name) == 0 {
		msg := "invalid definition name"
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return errors.New(msg)
	}

	if len(def.Description) == 0 {
		msg := "invalid description"
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return errors.New(msg)
	}

	//format SQL
	command := fmt.Sprintf("UPDATE %s SET name = ?, description = ? WHERE id = ?", table)

	//DO IT!!
	result, err := db.DB().Exec(command, def.Name, def.Description, def.ID)
	if err != nil {
		msg := fmt.Sprintf("unable to update designation: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	//make sure it acutally worked
	numRows, err := result.RowsAffected()
	if err != nil {
		msg := fmt.Sprintf("number of rows not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	if numRows < 1 {
		msg := "invalid edit"
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	return nil
}

func GetDefinitionById(table string, id int64, def *Definition) error {

	log.Printf("[accessors] fetching definition from %s with id %d", table, id)

	//format SQL
	command := fmt.Sprintf("SELECT * FROM %s WHERE id = ?", table)

	//check SQL
	log.Printf("SQL: %s", command)

	//fill struct
	err := db.DB().Get(def, command, id)
	if err != nil {
		msg := fmt.Sprintf("definition not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	return nil
}

func GetAllDefinitions(table string, defs *[]Definition) error {

	log.Printf("[accessors] getting all definitions from table: %s", table)

	cmd := fmt.Sprintf("SELECT * FROM %s", table)

	err := db.DB().Select(defs, cmd)
	if err != nil {
		msg := fmt.Sprintf("definitions not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	return nil
}
