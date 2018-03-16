package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/byuoitav/pi-designation-microservice/files"
	"github.com/labstack/echo"
)

//func ConvertYamlToBytes(microservices []ac.DBMicroservice) ([]byte, error) {
//
//	log.Printf("[handlers] converting microservice structs to text...")
//
//	var output bytes.Buffer
//
//	output.WriteString("version: '3.0'\nservices:\n") //common to all JSON
//
//	for _, microservice := range microservices {
//
//		output.WriteString(microservice.YAML)
//		output.WriteString("\n")
//	}
//
//	return output.Bytes(), nil
//
//}

func GetDockerComposeByDevice(context echo.Context) error {

	deviceId, err := strconv.Atoi(context.Param("id"))
	if err != nil {
		return context.JSON(http.StatusBadRequest, fmt.Sprintf("invalid device ID: %s", err.Error()))
	}

	hashCode, err := files.GetDockerComposeByDevice(deviceId)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	return context.JSON(http.StatusOK, hashCode)
}

func GetDockerComposeByRoomAndRole(context echo.Context) error {

	roomId, err := strconv.Atoi(context.Param("room"))
	if err != nil {
		return context.JSON(http.StatusBadRequest, "invalid room ID")
	}

	roleId, err := strconv.Atoi(context.Param("role"))
	if err != nil {
		return context.JSON(http.StatusBadRequest, "invalid role ID")
	}

	hashMap, err := files.GetDockerComposeByRoomAndRole(roomId, roleId)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	return context.JSON(http.StatusOK, hashMap)

}
