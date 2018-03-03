package files

import (
	"log"

	"github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/byuoitav/pi-designation-microservice/configuration"
	"github.com/fatih/color"
)

func GetEnvironmentByDevice(deviceId int) (string, error) {

	designationId, target, _, err := FetchDeviceMetaData(deviceId)
	if err != nil {
		return "", err
	}

	microservices, err := configuration.GetDockerComposeByDevice(target, designationId, map[int64]accessors.DBMicroservice{})
	if err != nil {
		return "", err
	}

	variables, err := configuration.GetDeviceEnvironment(microservices)
	if err != nil {
		return "", err
	}

	return WriteEnvironmentFile(variables)
}

func GetEnvironmentByRoomAndRole(roomId, roleId int) (map[int]string, error) {

	designationId, targets, _, err := FetchRoomMetaData(roomId, roleId)
	if err != nil {
		return nil, err
	}

	output := make(map[int]string)

	for _, target := range targets {

		log.Printf("%s", color.HiGreenString("considering target: %s", target.Name))

		microservices, err := configuration.GetDockerComposeByDevice(target, designationId, map[int64]accessors.DBMicroservice{})
		if err != nil {
			return nil, err
		}

		variables, err := configuration.GetDeviceEnvironment(microservices)
		if err != nil {
			return nil, err
		}

		fileName, err := WriteEnvironmentFile(variables)
		if err != nil {
			return nil, err
		}

		output[target.ID] = fileName

	}

	return output, nil
}
