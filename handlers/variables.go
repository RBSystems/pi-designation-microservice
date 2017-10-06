package handlers

import (
	"fmt"
	"log"
	"net/http"

	ac "github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/fatih/color"
	"github.com/labstack/echo"
)

func AddVariable(context echo.Context) error {

	//bind context
	var variable ac.Variable
	err := context.Bind(&variable)
	if err != nil {
		msg := fmt.Sprintf("unable to bind context: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	log.Printf("[handlers] handling request to add variable: %s", variable.Key)

	//validate key and value
	err = ac.ValidateVar(variable)
	if err != nil {
		msg := fmt.Sprintf("invalid variable: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	//validate designation
	variable.Desig, err = ac.GetDesignationByName(variable.Desig.Name)
	if err != nil {
		msg := fmt.Sprintf("invalid designation: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	//make sure it's not already there, this should error out because the variable isn't there
	//	exists := ac.FillVariable(&variable)
	//	if exists == nil {
	//		msg := fmt.Sprintf("variable: %s already present in database", variable.Key)
	//		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
	//		return context.JSON(http.StatusBadRequest, msg)
	//	}

	//add variable
	err = ac.AddNewVariable(variable)
	if err != nil {
		msg := fmt.Sprintf("unable to add new variable: %s: %s", variable.Key, err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusInternalServerError, msg)
	}

	return context.JSON(http.StatusOK, variable)
}

func GetVariable(context echo.Context) error {

	designation := context.Param("designation")
	key := context.Param("key")
	log.Printf("[handlers] getting %s value of %s", designation, key)

	//validate designation
	desig, err := ac.GetDesignationByName(designation)
	if err != nil {
		msg := fmt.Sprintf("%s designation not found: %s", designation, err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	//build variable
	variable := ac.Variable{
		Key:   key,
		Desig: desig,
	}

	//fill variable
	err = ac.FillVariable(&variable)
	if err != nil {
		msg := fmt.Sprintf("unable to get variable data: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusInternalServerError, msg)
	}

	return context.JSON(http.StatusOK, variable)

}

func EditVariable(context echo.Context) error {

	log.Printf("[handlers] updating variable: %s", context.Param("key"))

	//bind context
	var variable ac.Variable
	err := context.Bind(&variable)
	if err != nil {
		msg := fmt.Sprintf("unable to bind context: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	//validate
	err = ac.ValidateVar(variable)
	if err != nil {
		msg := fmt.Sprintf("invalid variable: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	//fill designation
	variable.Desig, err = ac.GetDesignationByName(variable.Desig.Name)
	if err != nil {
		msg := fmt.Sprintf("invalid designation: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	//make edit
	err = ac.EditVariable(variable)
	if err != nil {
		msg := fmt.Sprintf("unable to edit variable: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, variable)
}

func DeleteVariable(context echo.Context) error {

	key := context.Param("key")
	desig := context.Param("designation")
	log.Printf("[handlers] deleting %s-designated variable %s", key, desig)

	designation, err := ac.GetDesignationByName(desig)
	if err != nil {
		msg := fmt.Sprintf("invalid designation: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	variable := ac.Variable{
		Key:   key,
		Desig: designation,
	}

	err = ac.DeleteVariable(variable)
	if err != nil {
		msg := fmt.Sprintf("variable not deleted: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, "success")
}

func GetAllVariables(context echo.Context) error {

	log.Printf("[handlers] getting all variables...")

	vars, err := ac.GetAllVariables()
	if err != nil {
		msg := fmt.Sprintf("variables not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusInternalServerError, msg)
	}

	return context.JSON(http.StatusOK, vars)
}

func GetVarsByDesignation(context echo.Context) error {

	designation := context.Param("designation")
	log.Printf("[hanlders] getting all variables corresponding to designation: %s", designation)

	//build designation
	desig, err := ac.GetDesignationByName(designation)
	if err != nil {
		msg := fmt.Sprintf("designation %s not found: %s", designation, err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	//get vars by designation
	vars, err := ac.GetVariablesByDesignation(desig)
	if err != nil {
		msg := fmt.Sprintf("variables not found: %s", err.Error())
		log.Printf("%s", color.HiRedString("[handlers] %s", msg))
		return context.JSON(http.StatusBadRequest, msg)
	}

	return context.JSON(http.StatusOK, vars)
}
