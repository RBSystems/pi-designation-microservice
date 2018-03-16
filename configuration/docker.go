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

func GetDockerComposeByDevice(target structs.Device, designationId int64, commandMicroservices map[int64]accessors.Microservice) (map[int64]accessors.Microservice, error) {

	fmt.Printf("\t\t\t\tcommandMicroservices: ")
	for k, _ := range commandMicroservices {
		fmt.Printf(color.HiMagentaString("%d ", k))
	}
	fmt.Printf("\n")

	roles, err := GetTargetRoles(target) //	map role IDs from config DB to class IDs from desig DB
	if err != nil {
		msg := fmt.Sprintf("roles of %s not found: %s", target.Name, err.Error())
		log.Printf("%s", color.HiRedString("[configuration] %s", msg))
		return nil, errors.New(msg)
	}

	roleSet := make(map[int64]accessors.Microservice)      //	working minimum set of functionality microservices for device
	potentialSet := make(map[int64]accessors.Microservice) //	working potential set of command-oriented microservices

	for _, role := range roles {

		minimumSet, err := GetMinimumSet(role, designationId) //	there exists a minimum set for each role
		if err != nil {
			msg := fmt.Sprintf("minimum set for role ID %d not found: %s", role, err.Error())
			log.Printf("%s", color.HiRedString("[configuration] %s", msg))
			return nil, errors.New(msg)
		}

		roleSet = MicroserviceUnion(roleSet, minimumSet) //	the minimum set for a device is the union of all the roles

		possibleSet, err := GetPotentialSet(role, designationId) //	get the potential set of microservices for a device
		if err != nil {
			msg := fmt.Sprintf("potential microservices for %s not found: %s", target.Name, err.Error())
			log.Printf("%s", color.HiRedString("[configuration] %s", msg))
			return nil, errors.New(msg)
		}

		potentialSet = MicroserviceUnion(potentialSet, possibleSet)
	}

	fmt.Printf("\t\t\t\troleSet: ")
	for k, _ := range roleSet {
		fmt.Printf(color.HiMagentaString("%d ", k))
	}
	fmt.Printf("\n")

	fmt.Printf("\t\t\t\tpotentialSet: ")
	for k, _ := range potentialSet {
		fmt.Printf(color.HiMagentaString("%d ", k))
	}
	fmt.Printf("\n")

	commandSet := MicroserviceIntersect(potentialSet, commandMicroservices) //	find which microservices the device actually needs

	fmt.Printf("\t\t\t\tcommandSet: ")
	for k, _ := range commandSet {
		fmt.Printf(color.HiMagentaString("%d ", k))
	}
	fmt.Printf("\n")

	actualSet := MicroserviceUnion(commandSet, roleSet) //	a device needs the union of its role set and its command set

	fmt.Printf("\t\t\t\tfinalSet: ")
	for k, _ := range actualSet {
		fmt.Printf(color.HiMagentaString("%d ", k))
	}
	fmt.Printf("\n")

	return actualSet, nil // map the target's ID to the list of services
}

func GetRoomDesignationIdByRoomId(roomId int) (int64, error) {

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

func GetDeviceAndDesignationByDeviceId(deviceId int) (int64, structs.Device, error) {

	log.Printf("[configuration] finding designation of device %d", deviceId)

	device, err := dbo.GetDeviceById(deviceId)
	if err != nil {
		msg := fmt.Sprintf("device with id %d not found: %s", deviceId, err.Error())
		log.Printf("%s", color.HiRedString("[configuration] %s", msg))
		return 0, structs.Device{}, errors.New(msg)
	}

	log.Printf("[configuration] identified target device: %s", color.HiGreenString(device.Name))

	possibleDesignations, err := accessors.GetAllDesignations()
	if err != nil {
		msg := fmt.Sprintf("error fetching designations: %s", err.Error())
		log.Printf("%s", color.HiRedString("[configuration] %s", msg))
		return 0, structs.Device{}, errors.New(msg)
	}

	for _, designation := range possibleDesignations {

		if device.Room.RoomDesignation == designation.Name {

			log.Printf("[configuration] identified designation ID: %s", color.HiGreenString("%d (%s)", designation.ID, device.Room.RoomDesignation))

			return designation.ID, device, nil
		}
	}

	msg := fmt.Sprintf("designation of device: %s not found", device.Name)
	log.Printf("%s", color.HiRedString("[configuration] %s", msg))

	return 0, structs.Device{}, errors.New(msg)
}

//builds map of microservice_definition.id to yaml
func GetDeviceMicroservices(designationId int64, devices []structs.Device) (map[int64]accessors.Microservice, error) {

	output := make(map[int64]accessors.Microservice) //	map of microservice IDs to microservice mappings

	configToDesig, err := mapMicroservices()
	if err != nil {
		return nil, err
	}

	for _, device := range devices {

		log.Printf("considering device: %s", color.HiGreenString(device.Name))

		for _, command := range device.Commands {

			log.Printf("\tconsidering command: %s (%s)", color.HiBlueString(command.Name), color.HiBlueString(command.Microservice))

			microserviceId := configToDesig[command.Microservice]

			if !strings.Contains(command.Microservice, "gateway") { //	TODO come up with a better way to handle this

				if _, ok := output[microserviceId]; !ok {

					microserviceMapping, err := accessors.GetMicroserviceMappingByDesignation(designationId, microserviceId)
					if err != nil {
						return nil, err
					}

					output[microserviceMapping.MicroserviceId] = microserviceMapping

					log.Printf("%s", color.HiYellowString("\t\t added %d", microserviceId))
				}
			}
		}
	}

	return output, nil
}

func GetMinimumSet(class, designation int64) (map[int64]accessors.Microservice, error) {

	var set []accessors.Microservice
	err := accessors.GetMinimumSet(&set, class, designation)
	if err != nil {
		return nil, err
	}

	output := make(map[int64]accessors.Microservice)

	for _, microservice := range set {

		output[microservice.MicroserviceId] = microservice //	FIXME IDs come from standard_sets table (they shouldn't)
	}

	return output, nil
}

func GetPotentialSet(class, designation int64) (map[int64]accessors.Microservice, error) {

	var set []accessors.Microservice
	err := accessors.GetPossibleSet(&set, class, designation)
	if err != nil {
		return nil, err
	}

	output := make(map[int64]accessors.Microservice)

	for _, microservice := range set {

		output[microservice.MicroserviceId] = microservice //	FIXME IDs come from wrong table
	}

	return output, nil
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

	return output, nil
}

func MicroserviceUnion(a map[int64]accessors.Microservice, b map[int64]accessors.Microservice) map[int64]accessors.Microservice {

	for key, value := range a {

		if _, ok := b[key]; !ok {

			b[key] = value
		}
	}

	return b
}

func MicroserviceIntersect(a map[int64]accessors.Microservice, b map[int64]accessors.Microservice) map[int64]accessors.Microservice {

	intersect := make(map[int64]accessors.Microservice)

	for key, value := range a { // consider each element in a

		if _, ok := b[key]; ok { // if that element is also in b, the element is a member of the intersect

			intersect[key] = value
		}
	}

	return intersect
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

				output[cMicro.Address] = dMicro.Id
			}
		}
	}

	return output, nil
}
