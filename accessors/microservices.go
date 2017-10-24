package accessors

import (
	"errors"
	"log"

	db "github.com/byuoitav/pi-designation-microservice/database"
	"github.com/fatih/color"
)

func AddMicroserviceDefinition(microservice *MicroserviceDefinition) error {

	if len(microservice.Name) == 0 {
		msg := "invalid microservice name"
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	log.Printf("[accessors] adding new microservice: %s", microservice.Name)

	microDef := "INSERT INTO microservice_definitions (name, description) VALUES (?, ?)"
	result, err := db.DB().Exec(microDef, microservice.Name, microservice.Description)
	if err != nil {
		log.Printf("%s", color.HiRedString("[accessors] %s", err.Error()))
		return err
	}

	microservice.ID, err = result.LastInsertId()
	if err != nil {
		log.Printf("%s", color.HiRedString("[accessors] %s", err.Error()))
		return err
	}

	return nil
}
