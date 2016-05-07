package utils

import "log"

// CheckError logs error
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
