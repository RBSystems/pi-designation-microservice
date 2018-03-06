package accessors

import (
	"log"

	db "github.com/byuoitav/pi-designation-microservice/database"
)

type Microservice struct {
	Id             int64  `db:"id"` //FIXME do we really need this?
	MicroserviceId int64  `db:"microservice_id"`
	Name           string `db:"name"`
	DesignationId  int64  `db:"designation_id"`
	ClassId        int64  `db:"class_id"`
	Yaml           string `db:"yaml"`
}

type MicroserviceDefinition struct {
	Id          int64  `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}

func GetAllMicroservices() ([]MicroserviceDefinition, error) {

	var output []MicroserviceDefinition
	err := db.DB().Select(&output, "SELECT * FROM microservice_definitions")
	if err != nil {
		return []MicroserviceDefinition{}, err
	}

	return output, nil
}

func GetMicroserviceMappingByDesignation(designationId, microserviceId int64) (Microservice, error) {

	log.Printf("[accessors] getting microservice mapping with designation ID: %d and microservice ID: %d", designationId, microserviceId)

	var output Microservice

	command := "SELECT * FROM microservice_mappings WHERE designation_id = ? AND microservice_id = ?"

	err := db.DB().Get(&output, command, designationId, microserviceId)
	if err != nil {
		return Microservice{}, err
	}

	return output, nil
}
