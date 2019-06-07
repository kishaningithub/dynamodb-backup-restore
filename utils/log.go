package utils

import "log"

func CheckError(message string, err error) {
	if err != nil {
		log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
		log.Fatal(message, ": ", err.Error())
	}
}

func PrintInfo(message string) {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.Println(message)
}
