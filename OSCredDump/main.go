package main

import (
	"helpers"
	"os"
	"runtime"
	"strings"
)

func main() {
	nameOfLogFile := helpers.CreateLogFileIfItDoesNotExist("./", "OSDump")
	helpers.WriteLog(nameOfLogFile, "Starting test: OsCredDump", 2)
	currOS := runtime.GOOS
	if currOS == "windows" {
		paths := os.Getenv("SystemRoot")
		listOfPaths := strings.Split(paths, string(os.PathListSeparator))

		_, err := os.OpenFile(listOfPaths[0]+"\\System32\\config\\SAM", os.O_RDWR, 0666)
		if err != nil {

			helpers.WriteLog(nameOfLogFile, err.Error(), 1)
			os.Exit(1)
		} else {
			helpers.WriteLog(nameOfLogFile, "Ending test: OsCredDump", 2)
			os.Exit(0)
		}
	}
}
