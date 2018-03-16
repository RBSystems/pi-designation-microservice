package configuration

import (
	"fmt"
	"log"

	"github.com/byuoitav/configuration-database-microservice/structs"
	"github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/fatih/color"
)

const PI_HOSTNAME = 64

func GetDeviceEnvironment(target *structs.Device, microservices map[int64]accessors.Microservice) (map[int64]accessors.Variable, error) {

	output := make(map[int64]accessors.Variable)

	for _, microservice := range microservices {

		log.Printf("[configuration] considering microservice: %s", color.HiMagentaString(microservice.Name))

		potentialVars, err := accessors.GetVariablesByMicroserviceAndDesignation(microservice.MicroserviceId, microservice.DesignationId)
		if err != nil {
			return nil, err
		}

		fmt.Printf("potential set: \n")
		for _, v := range potentialVars {
			fmt.Printf("\t\t%s\n", color.HiMagentaString(v.Name))
		}
		fmt.Printf("\n")

		output = VariableUnion(potentialVars, output)

		fmt.Printf("working set: \n")
		for _, v := range output {
			fmt.Printf("\t\t%s\n", color.HiMagentaString(v.Name))
		}
		fmt.Printf("\n")
	}

	hostname := fmt.Sprintf("%s-%s-%s", target.Building.Shortname, target.Room.Name, target.Name) // add discrete hostname
	output[PI_HOSTNAME] = accessors.Variable{Name: "PI_HOSTNAME", Value: hostname}

	return output, nil
}

func VariableUnion(a, b map[int64]accessors.Variable) map[int64]accessors.Variable {

	for k, v := range a {

		b[k] = v
	}

	return b
}

func VariableIntersect(a, b map[int64]accessors.Variable) map[int64]accessors.Variable {

	intersect := make(map[int64]accessors.Variable)

	for k, v := range a {

		if _, ok := b[k]; ok {

			intersect[k] = v
		}
	}

	return intersect
}

func VariableContains(a []accessors.Variable, v accessors.Variable) bool {

	for _, variable := range a {

		if VariableEquals(variable, v) {

			return true
		}

	}

	return false
}

func VariableEquals(a, b accessors.Variable) bool {

	if a.Id != b.Id ||
		a.DesignationId != b.DesignationId ||
		a.DefinitionId != b.DefinitionId ||
		a.MicroserviceId != b.MicroserviceId ||
		a.Name != b.Name ||
		a.Value != b.Value {
		return false
	}

	return true
}
