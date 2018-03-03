package files

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/byuoitav/pi-designation-microservice/accessors"
	"github.com/fatih/color"
)

var TIMER_DURATION = 3 * time.Minute
var fileTimers sync.Map

func GenerateRandomString(numBytes int) (string, error) {

	bytes := make([]byte, numBytes)
	if _, err := rand.Read(bytes); err != nil {
		return "", errors.New(fmt.Sprintf("error generating file name: %s", err.Error()))
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

func TrackFile(path string) {
	removeFile := func() {
		log.Printf("[helpers] removing old file: %s...", color.HiYellowString(path))
		err := os.Remove(path)
		if err != nil {
			log.Printf("%s", color.HiRedString("[helpers] error removing old file: %s", err.Error()))
		}
	}

	fileTimers.Store(path, time.AfterFunc(TIMER_DURATION, removeFile))
}

// calls Hash and TrackFile
func WriteDockerComposeFile(microservices map[int64]accessors.DBMicroservice) (string, error) {
	fileLocation := fmt.Sprintf("%s/src/github.com/byuoitav/pi-designation-microservice/public", os.Getenv("GOPATH"))
	fileName, err := GenerateRandomString(len(microservices))
	if err != nil {
		return "", err
	}

	path := fmt.Sprintf("%s/%s", fileLocation, fileName)

	_, err = os.Stat(path)
	for err == nil {
		fileName, err = GenerateRandomString(len(microservices))

		path = fmt.Sprintf("%s/%s", fileLocation, fileName)
	}

	outFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		return "", err
	}

	outFile.WriteString("version: '3'\nservices:\n") // common to all YAML files we want

	for _, microservice := range microservices {

		outFile.WriteString(microservice.YAML)
		outFile.WriteString("\n")
	}

	outFile.Close()

	TrackFile(path)

	return fileName, nil
}

func WriteEnvironmentFile(variables map[int64]accessors.Variable) (string, error) {
	return "", nil
}
