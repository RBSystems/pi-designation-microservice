package accessors

import (
	"log"

	db "github.com/byuoitav/pi-designation-microservice/database"
)

func GetMicroserviceMappingByDesignation(designationId, microserviceId int64) (DBMicroservice, error) {

	log.Printf("[accessors] getting microservice mapping with designation ID: %d and microservice ID: %d", designationId, microserviceId)

	var output DBMicroservice

	command := "SELECT * FROM microservice_mappings WHERE designation_id = ? AND microservice_id = ?"

	err := db.DB().Get(&output, command, designationId, microserviceId)
	if err != nil {
		return DBMicroservice{}, err
	}

	return output, nil
}
