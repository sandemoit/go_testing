package utils

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)

func LogError(err error) {
	Logger.Fatal(err)
}

func LogInfo(msg string) {
	Logger.Println(msg)
}
