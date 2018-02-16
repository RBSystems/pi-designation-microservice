package configuration

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/byuoitav/av-api/dbo"
	"github.com/byuoitav/configuration-database-microservice/structs"
	"github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/fatih/color"
)

/*
building a set of microservice definitions, then query with the designation ID as a parameter
*/

func GetDockerComposeByRoomAndRole(roomId, roleId int) (map[int][]accessors.DBMicroservice, error) {

	devices, err := dbo.GetDevicesByRoomId(roomId) //	get all devices in room
	if err != nil {
		msg := fmt.Sprintf("devices not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[configuration] %s", msg))
		return nil, errors.New(msg)
	}

	for _, device := range devices {
		log.Printf("%s", color.HiBlueString("\t%s", device.Name))
	}

	designationId, err := GetRoomDesignationId(roomId)
	if err != nil {
		msg := fmt.Sprintf("room designation not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[configuration] %s", msg))
		return nil, errors.New(msg)
	}

	commandMicroservices, err := GetDeviceMicroservices(designationId, devices) // build set of all microservices the devices might need
	if err != nil {
		msg := fmt.Sprintf("device command microservices not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[configuration] %s", msg))
		return nil, errors.New(msg)
	}

	targets, err := dbo.GetDevicesByRoomIdAndRoleId(roomId, roleId) //	get target devices
	if err != nil {
		msg := fmt.Sprintf("targets not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[configuration] %s", msg))
		return nil, errors.New(msg)
	}

	output := make(map[int][]accessors.DBMicroservice) //	maps a device ID to a list of YAML snippets

	for _, target := range targets { // build a list of microservices for each target

		log.Printf("%s", color.HiGreenString("considering target: %s", target.Name))

		roles, err := GetTargetRoles(target) //	map role IDs from config DB to class IDs from desig DB
		if err != nil {
			msg := fmt.Sprintf("roles of %s not found: %s", target.Name, err.Error())
			log.Printf("%s", color.HiRedString("[configuration] %s", msg))
			return nil, errors.New(msg)
		}

		var roleSet map[int64]accessors.DBMicroservice //	working minimum set of microservices for device
		for _, role := range roles {

			minimumSet, err := GetMinimumSet(role, designationId) //	there exists a minimum set for each role
			if err != nil {
				msg := fmt.Sprintf("minimum set for role ID %d not found: %s", role, err.Error())
				log.Printf("%s", color.HiRedString("[configuration] %s", msg))
				return nil, errors.New(msg)
			}

			log.Printf("%v", minimumSet)

			roleSet = MicroserviceUnion(roleSet, minimumSet) //	the minimum set for a device is the union of all the roles
		}

		log.Printf("%v", roleSet)

		potentialSet, err := GetPotentialSet(target) //	get the potential set of microservices for a device
		if err != nil {
			msg := fmt.Sprintf("potential microservices for %s not found: %s", target.Name, err.Error())
			log.Printf("%s", color.HiRedString("[configuration] %s", msg))
			return nil, errors.New(msg)
		}

		commandSet := MicroserviceIntersect(potentialSet, commandMicroservices) //	find which microservices the device actually needs
		actualSet := MicroserviceUnion(commandSet, roleSet)                     //	a device needs the union of its role set and its command set

		output[target.ID] = convertToList(actualSet) // map the target's ID to the list of services
	}

	return nil, nil
}

func GetRoomDesignationId(roomId int) (int64, error) {

	room, err := dbo.GetRoomById(roomId)
	if err != nil {
		return 0, err
	}

	possibleDesignations, err := accessors.GetAllDesignations()
	if err != nil {
		return 0, err
	}

	for _, possibleDesignation := range possibleDesignations {

		if room.RoomDesignation == possibleDesignation.Name {

			log.Printf("%s", color.HiGreenString("identified designation ID: %d (%s)", possibleDesignation.ID, room.RoomDesignation))
			return possibleDesignation.ID, nil
		}
	}

	return 0, errors.New(fmt.Sprintf("%d is not a valid designation ID"))
}

//builds map of microservice_definition.id to yaml
func GetDeviceMicroservices(designationId int64, devices []structs.Device) (map[int64]accessors.DBMicroservice, error) {

	output := make(map[int64]accessors.DBMicroservice) //	map of microservice IDs to microservice mappings

	configToDesig, err := mapMicroservices()
	if err != nil {
		return nil, err
	}

	for _, device := range devices {

		log.Printf("considering device: %s", color.HiGreenString(device.Name))

		for _, command := range device.Commands {

			log.Printf("\tconsidering command: %s (%s)", color.HiBlueString(command.Name), color.HiBlueString(command.Microservice))

			microserviceId := configToDesig[command.Microservice]

			if _, ok := output[microserviceId]; !ok {

				microserviceMapping, err := accessors.GetMicroserviceMappingByDesignation(designationId, microserviceId)
				if err != nil {
					return nil, err
				}

				output[microserviceId] = microserviceMapping

				log.Printf("%s", color.HiYellowString("\t\t added %d", microserviceId))
			}
		}
	}

	return map[int64]accessors.DBMicroservice{}, nil
}

func GetMinimumSet(class, designation int64) (map[int64]accessors.DBMicroservice, error) {

	var set []accessors.DBMicroservice
	err := accessors.GetMinimumSet(&set, class, designation)
	if err != nil {
		return nil, err
	}

	output := make(map[int64]accessors.DBMicroservice)

	for _, microservice := range set {

		output[microservice.ID] = microservice
	}

	return output, nil
}

func GetPotentialSet(target structs.Device) (map[int64]accessors.DBMicroservice, error) {

	return nil, nil
}

func GetTargetRoles(target structs.Device) ([]int64, error) {

	allClasses, err := accessors.GetAllClasses()
	if err != nil {
		return []int64{}, err
	}

	var output []int64
	for _, role := range target.Roles {

		for _, class := range allClasses {

			if class.Name == role {

				log.Printf("%s", color.HiCyanString("\tfound role %s", role))
				output = append(output, class.ID)
			}
		}
	}

	return []int64{}, nil
}

func MicroserviceUnion(a map[int64]accessors.DBMicroservice, b map[int64]accessors.DBMicroservice) map[int64]accessors.DBMicroservice {

	for key, value := range a {

		if _, ok := b[key]; !ok {

			b[key] = value
		}
	}

	return b
}

func MicroserviceIntersect(a map[int64]accessors.DBMicroservice, b map[int64]accessors.DBMicroservice) map[int64]accessors.DBMicroservice {

	return nil
}

func convertToList(a map[int64]accessors.DBMicroservice) []accessors.DBMicroservice {

	return []accessors.DBMicroservice{}
}

func getMicroserviceName(address string) (string, error) {

	return "", errors.New("microservice name not found")
}

func mapMicroservices() (map[string]int64, error) {

	desigMicros, err := accessors.GetAllMicroservices() //	find the set of microservices the designation db knows about
	if err != nil {
		return nil, err
	}

	configMicros, err := dbo.GetMicroservices()
	if err != nil {
		return nil, err
	}

	output := make(map[string]int64)

	for _, cMicro := range configMicros {

		for _, dMicro := range desigMicros {

			if strings.Compare(cMicro.Name, dMicro.Name) == 0 {

				log.Printf("mapping %s to %d", cMicro.Address, dMicro.ID)
				output[cMicro.Address] = dMicro.ID
			}
		}
	}

	log.Printf("output: %d, configMicros: %d", len(output), len(configMicros))
	log.Printf("%v", output)

	if len(output) < len(configMicros) {
		log.Printf("%s", color.HiRedString("[configuration] warning: some configuration microservices did not map to designation microservices!"))
	}

	return output, nil
}
