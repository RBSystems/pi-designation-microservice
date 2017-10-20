package accessors

import (
	"errors"
	"fmt"
	"log"

	"github.com/byuoitav/pi-designation-microservice/database"
	"github.com/fatih/color"
)

func AddMicroserviceDefinition(microservice *MicroserviceDefinition) error {

	log.Printf("[accessors] adding %s definition to database...", microservice.Name)

	_, err := database.DB().Exec("INSERT into microservice definitions (name, yaml, description) values(?, ?, ?)", microservice.Name, microservice.YAML, microservice.Description)
	if err != nil {
		msg := fmt.Sprintf("unable to add microservice: %s", err.Error())
		log.Printf("%s", color.HiRedString("[accessors] %s", msg))
		return errors.New(msg)
	}

	return nil
}
