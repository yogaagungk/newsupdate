package common

import "log"

// HandleError to handle error message on log
func HandleError(err error, message string) {
	if err != nil {
		log.Fatal(message, err)
	}
}
