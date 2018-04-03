package files

import (
	"errors"
	"fmt"
	"log"

	"github.com/byuoitav/av-api/dbo"
	"github.com/byuoitav/configuration-database-microservice/structs"
	"github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/byuoitav/pi-designation-microservice/configuration"
	"github.com/fatih/color"
)

func GetDockerComposeByDevice(deviceId int) (string, error) {

	log.Printf("%s", color.HiGreenString("[files] fetching docker-compose for %d", deviceId))

	designationId, device, microservices, err := FetchDeviceMetaData(deviceId)
	if err != nil {
		msg := fmt.Sprintf("device meta-data not found: %s", err.Error())
		return "", errors.New(msg)
	}

	deviceSet, err := configuration.GetDockerComposeByDevice(device, designationId, microservices)
	if err != nil {
		msg := fmt.Sprintf("device microservices not found: %s", err.Error())
		return "", errors.New(msg)
	}

	return WriteDockerComposeFile(deviceSet)
}

func GetDockerComposeByRoomAndRole(roomId int, roleId int) (map[int]string, error) {

	//	get room metadata
	designationId, targets, commandMicroservices, err := FetchRoomMetaData(roomId, roleId)
	if err != nil {
		return nil, err
	}

	output := make(map[int]string) // build a map of device IDs to hash codes

	for _, target := range targets { // build a list of microservices for each target

		log.Printf("%s", color.HiGreenString("considering target: %s", target.Name))

		deviceSet, err := configuration.GetDockerComposeByDevice(target, designationId, commandMicroservices)
		if err != nil {
			return nil, err
		}

		fileName, err := WriteDockerComposeFile(deviceSet)
		if err != nil {
			return nil, err
		}

		output[target.ID] = fileName

	}

	return output, nil
}

func FetchDeviceMetaData(deviceId int) (int64, structs.Device, map[int64]accessors.Microservice, error) {

	log.Printf("[files] fetching meta-data for %d", deviceId)

	designationId, device, err := configuration.GetDeviceAndDesignationByDeviceId(deviceId)
	if err != nil {
		return 0, structs.Device{}, nil, err
	}

	devices, err := dbo.GetDevicesByRoomId(device.Room.ID)
	if err != nil {
		return 0, structs.Device{}, nil, err
	}

	commandMicroservices, err := configuration.GetDeviceMicroservices(designationId, devices)
	if err != nil {
		return 0, structs.Device{}, nil, err
	}

	return designationId, device, commandMicroservices, nil

}

func FetchRoomMetaData(roomId int, roleId int) (int64, []structs.Device, map[int64]accessors.Microservice, error) {

	devices, err := dbo.GetDevicesByRoomId(roomId) //	get all devices in room
	if err != nil {
		msg := fmt.Sprintf("devices not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[configuration] %s", msg))
		return 0, nil, nil, errors.New(msg)
	}

	for _, device := range devices {
		log.Printf("%s", color.HiBlueString("\t%s", device.Name))
	}

	designationId, err := configuration.GetRoomDesignationIdByRoomId(roomId)
	if err != nil {
		msg := fmt.Sprintf("room designation not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[configuration] %s", msg))
		return 0, nil, nil, errors.New(msg)
	}

	commandMicroservices, err := configuration.GetDeviceMicroservices(designationId, devices) // build set of all microservices the devices might need
	if err != nil {
		msg := fmt.Sprintf("device command microservices not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[configuration] %s", msg))
		return 0, nil, nil, errors.New(msg)
	}

	targets, err := dbo.GetDevicesByRoomIdAndRoleId(roomId, roleId) //	get target devices
	if err != nil {
		msg := fmt.Sprintf("targets not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[configuration] %s", msg))
		return 0, nil, nil, errors.New(msg)
	}

	fmt.Printf("\t\t\t\ttargets: ")
	for _, t := range targets {
		fmt.Printf(color.HiMagentaString(t.Name))
	}
	fmt.Printf("\n")

	fmt.Printf("\t\t\t\tcommandMicroservices: ")
	for k, _ := range commandMicroservices {
		fmt.Printf(color.HiMagentaString("%d ", k))
	}
	fmt.Printf("\n")

	return designationId, targets, commandMicroservices, nil

}
