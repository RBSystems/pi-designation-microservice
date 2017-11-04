package handlers

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/labstack/echo"
)

func ExtractId(context echo.Context) (int64, error) {

	stringId := context.Param("id")
	intId, err := strconv.Atoi(stringId)
	if err != nil {
		msg := fmt.Sprintf("invalid ID: %s", err.Error())
		return 0, errors.New(msg)
	}

	return int64(intId), nil

}
