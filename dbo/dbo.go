package dbo

import "log"

type Device struct {
	Hostname    string            `json:"name"`
	Designation string            `json:"designation"`
	Environment map[string]string `json:"environment"`
}

func GetEnv(hostname string) (Device, error) {

	log.Printf("[dbo] getting environment variables of hostname: %s", hostname)

	return Device{}, nil
}

func GetUi(hostname string) (string, error) {

	log.Printf("[dbo] getting ui configuration of %s", hostname)

	return "", nil
}
