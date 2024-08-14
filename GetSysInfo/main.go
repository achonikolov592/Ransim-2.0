package main

import (
	getinfo "RRA/GetInfo"
	"helpers"
	"os"
)

func main() {
	var nameOfLogFile string
	if len(os.Args) == 2 {
		nameOfLogFile = os.Args[1]
	} else {
		nameOfLogFile = helpers.CreateLogFileIfItDoesNotExist("./", "GetSystemInformation", "GetSystemInformation")
	}

	helpers.WriteLog(nameOfLogFile, "Starting Test: GetSysInfo", 2, "GetSystemInformation")
	os.Remove("./SystemInformation.log")
	nameOfContentFile := helpers.CreateNormalLogFileIfItDoesNotExist("./", "SystemInformation")

	err := getinfo.GetSysInfo(nameOfContentFile)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "GetSystemInformation")
		os.Exit(1)
	}
	helpers.WriteLog(nameOfLogFile, "Ending Test: GetSysInfo", 2, "GetSystemInformation")
	os.Exit(0)

}
