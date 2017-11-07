package handlers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/byuoitav/pi-designation-microservice/accessors"
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

	vars, err := accessors.GetVariablesByClassAndDesignation(int64(classInt), int64(desigInt))
	if err != nil {
		msg := fmt.Sprintf("variables not found: %s", err.Error())
		log.Printf("%s", color.HiRedString)
		return context.JSON(http.StatusBadRequest, msg)
	}

	file, err := ConvertVariablesToString(vars)
	if err != nil {
		msg := fmt.Sprintf("error converting variables to text: %s", err.Error())
		log.Printf("%s", color.HiRedString)
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.Blob(http.StatusOK, "text/plain", file)
}

func ConvertVariablesToString(vars []accessors.VariableMapping) ([]byte, error) {

	log.Printf("[handlers] converting variable structs to text...")
	var output bytes.Buffer

	for _, variable := range vars {

		output.WriteString(variable.Variable.Name)
		output.WriteString(" ")
		output.WriteString(variable.Value)
		output.WriteString("\n")
	}

	return output.Bytes(), nil
}
