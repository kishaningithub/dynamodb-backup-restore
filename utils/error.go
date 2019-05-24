package utils

import "log"

func CheckError(message string, err error) {
	if err != nil {
		log.Fatal(message, ": ", err.Error())
	}
}
