package common

import "log"

func HandleError(err error, message string) {
	if err != nil {
		log.Fatal("%s: %s", message, err)
	}
}

type StatusCode struct {
	Status string
}
