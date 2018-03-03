package configuration

import (
	"fmt"

	"github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/fatih/color"
)

func GetDeviceEnvironment(microservices map[int64]accessors.DBMicroservice) (map[int64]accessors.Variable, error) {

	output := make(map[int64]accessors.Variable)

	for _, microservice := range microservices {

		potentialVars, err := accessors.GetVariablesByMicroserviceAndDesignation(microservice.ID, microservice.DesigID)
		if err != nil {
			return nil, err
		}

		fmt.Printf("potential set: ")
		for _, v := range potentialVars {
			fmt.Printf("%s ", color.HiMagentaString(v.Name))
		}
		fmt.Printf("\n")

		//output = VariableIntersect(potentialVars, output)

		fmt.Printf("working set: ")
		for _, v := range output {
			fmt.Printf("%s ", color.HiMagentaString(v.Name))
		}
		fmt.Printf("\n")
	}

	return output, nil
}

func VariableIntersect(a, b []accessors.Variable) []accessors.Variable {

	var output []accessors.Variable

	for _, v := range a { //	 range over a, adding anyting b doesn't already have to b

		if !VariableContains(b, v) {

			b = append(b, v)
		}
	}

	return output
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
