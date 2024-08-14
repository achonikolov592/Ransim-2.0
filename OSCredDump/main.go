package main

import (
	"helpers"
	"os"
	"runtime"
	"strings"
)

func main() {
	var nameOfLogFile string
	if len(os.Args) == 2 {
		nameOfLogFile = os.Args[1]
	} else {
		nameOfLogFile = helpers.CreateLogFileIfItDoesNotExist("./", "OSCredentialDump", "OSCredentialDump")
	}

	helpers.WriteLog(nameOfLogFile, "Starting test: OsCredDump", 2, "OSCredentialDump")
	currOS := runtime.GOOS
	if currOS == "windows" {
		paths := os.Getenv("SystemRoot")
		listOfPaths := strings.Split(paths, string(os.PathListSeparator))

		_, err := os.OpenFile(listOfPaths[0]+"\\System32\\config\\SAM", os.O_RDWR, 0666)
		if err != nil {
			if strings.Contains(err.Error(), "Access is denied.") {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1, "OSCredentialDump")
				os.Exit(1)
			} else {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1, "OSCredentialDump")
				os.Exit(2)
			}
		}
		helpers.WriteLog(nameOfLogFile, "Ending test: OsCredDump", 2, "OSCredentialDump")
		os.Exit(0)
	}
}
