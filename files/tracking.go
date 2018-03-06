package files

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/fatih/color"
)

const RANDOM_LENGTH = 10

var TIMER_DURATION = 3 * time.Minute
var fileTimers sync.Map

func GenerateRandomString(numBytes int) (string, error) {

	bytes := make([]byte, numBytes)
	if _, err := rand.Read(bytes); err != nil {
		return "", errors.New(fmt.Sprintf("error generating file name: %s", err.Error()))
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

func WriteFile(content string) (string, error) {

	fileLocation := fmt.Sprintf("%s/src/github.com/byuoitav/pi-designation-microservice/public", os.Getenv("GOPATH"))
	fileName, err := GenerateRandomString(RANDOM_LENGTH)
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("%s/%s", fileLocation, fileName)

	_, err = os.Stat(path)
	for err == nil { // make sure we don't have a duplicate file name

		fileName, err = GenerateRandomString(RANDOM_LENGTH)
		path = fmt.Sprintf("%s/%s", fileLocation, fileName)
		_, err = os.Stat(path)
	}

	outFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return "", err
	}

	if _, err = outFile.WriteString(content); err != nil {
		return "", nil
	}

	outFile.Close()

	removeFile := func() {
		log.Printf("[helpers] removing old file: %s...", color.HiYellowString(path))
		err := os.Remove(path)
		if err != nil {
			log.Printf("%s", color.HiRedString("[helpers] error removing old file: %s", err.Error()))
		}
	}

	fileTimers.Store(path, time.AfterFunc(TIMER_DURATION, removeFile))

	return fileName, nil
}

// calls Hash and TrackFile
func WriteDockerComposeFile(microservices map[int64]accessors.Microservice) (string, error) {

	var file strings.Builder

	file.WriteString("version: '3'\nservices:\n") // common to all YAML files we want

	for _, microservice := range microservices {

		file.WriteString(fmt.Sprintf("%s\n", microservice.Yaml))
	}

	return WriteFile(file.String())

}

func WriteEnvironmentFile(variables map[int64]accessors.Variable) (string, error) {

	var file strings.Builder

	for _, variable := range variables {

		file.WriteString(fmt.Sprintf("%s='%s'\n", variable.Name, variable.Value))
	}

	return WriteFile(file.String())
}
