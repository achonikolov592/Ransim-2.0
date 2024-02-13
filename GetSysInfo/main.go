package main

import (
	getinfo "RRA/GetInfo"
	"helpers"
	"os"
)

func main() {
	nameOfLogFile := helpers.CreateLogFileIfItDoesNotExist("./", "GetSysInfo")
	helpers.WriteLog(nameOfLogFile, "Starting Test: GetSysInfo", 2)
	os.Remove("./SystemInformation.log")
	nameOfContentFile := helpers.CreateLogFileIfItDoesNotExist("./", "SystemInformation")

	err := getinfo.GetSysInfo(nameOfContentFile, nameOfLogFile)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(1)
	}
	helpers.WriteLog(nameOfLogFile, "Ending Test: GetSysInfo", 2)
	os.Exit(0)

}
