package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"

	ac "github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/byuoitav/pi-designation-microservice/configuration"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

func GetVariablesByDesignationAndClass(context echo.Context) error {

	desig := context.Param("designation")
	desigInt, err := strconv.Atoi(desig)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	class := context.Param("class")
	classInt, err := strconv.Atoi(class)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	log.Printf("%s", color.HiCyanString("[handlers] fetching all variables from desigation: %d, class: %d", desigInt, classInt))

	vars, err := ac.GetVariablesByClassAndDesignation(int64(classInt), int64(desigInt))
	if err != nil {
		msg := fmt.Sprintf("variables not found: %s", err.Error())
		log.Printf("%s", color.HiRedString)
		return context.JSON(http.StatusBadRequest, msg)
	}

	file, err := ConvertVariablesToBytes(vars)
	if err != nil {
		msg := fmt.Sprintf("error converting variables to text: %s", err.Error())
		log.Printf("%s", color.HiRedString)
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.Blob(http.StatusOK, "text/plain", file)
}

func ConvertVariablesToBytes(vars []ac.VariableMapping) ([]byte, error) {

	log.Printf("[handlers] converting variable structs to text...")
	var output bytes.Buffer

	for _, variable := range vars {

		output.WriteString(variable.Variable.Name)
		output.WriteString("=")
		output.WriteString(variable.Value)
		output.WriteString("\n")
	}

	return output.Bytes(), nil
}

func GetDockerComposeByDesignationAndClass(context echo.Context) error {

	desig := context.Param("designation")
	desigInt, err := strconv.Atoi(desig)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	class := context.Param("class")
	classInt, err := strconv.Atoi(class)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	log.Printf("%s", color.HiCyanString("[handlers] fetching all variables from desigation: %d, class: %d", desigInt, classInt))

	var yamlSnippets []ac.DBMicroservice
	err = ac.GetDockerComposeByDesignationAndClass(&yamlSnippets, int64(classInt), int64(desigInt))
	if err != nil {
		msg := fmt.Sprintf("docker-compose data not found: %s", err.Error())
		log.Printf("%s", color.HiRedString)
		return context.JSON(http.StatusBadRequest, msg)
	}

	file, err := ConvertYamlToBytes(yamlSnippets)
	if err != nil {
		msg := fmt.Sprintf("unable to parse YAML: %s", err.Error())
		log.Printf("%s", color.HiRedString)
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.Blob(http.StatusOK, "text/plain", file)
}

func ConvertYamlToBytes(microservices []ac.DBMicroservice) ([]byte, error) {

	log.Printf("[handlers] converting microservice structs to text...")

	var output bytes.Buffer

	output.WriteString("version: '3'\nservices:\n") //common to all JSON

	for _, microservice := range microservices {

		output.WriteString(microservice.YAML)
		output.WriteString("\n")
	}

	return output.Bytes(), nil

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

	yamlSnippets, err := configuration.GetDockerComposeByRoomAndRole(roomId, roleId)
	if err != nil {
		return context.JSON(http.StatusBadRequest, err.Error())
	}

	dockerCompose, err := ConvertYamlToBytes(yamlSnippets)
	if err != nil {
		return context.JSON(http.StatusInternalServerError, err.Error())
	}

	return context.Blob(http.StatusOK, "text/plain", dockerCompose)

}
