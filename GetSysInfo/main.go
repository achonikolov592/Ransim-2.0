package main

import (
	"helpers"
	"os"
	"os/exec"
)

func main() {
	nameOfLogFile := helpers.CreateLogFileIfItDoesNotExist("./", "GetSysInfo")
	helpers.WriteLog(nameOfLogFile, "Starting Test: GetSysInfo", 2)
	os.Remove("./SystemInformation.log")
	nameOfContentFile := helpers.CreateLogFileIfItDoesNotExist("./", "SystemInformation")
	getSystemInfo := exec.Command("systeminfo")
	getPowershellSystemInfo := exec.Command("powershell.exe", "Get-ComputerInfo")
	getTasklist := exec.Command("tasklist")
	getNetstat := exec.Command("netstat", "-a")
	getPortocolStatistics := exec.Command("netstat", "-s")
	getIpconfig := exec.Command("ipconfig", "/all")

	f, err := os.OpenFile(nameOfContentFile, os.O_APPEND, 0666)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(2)
	}

	getSystemInfo.Stdout = f
	getPowershellSystemInfo.Stdout = f
	getTasklist.Stdout = f
	getNetstat.Stdout = f
	getPortocolStatistics.Stdout = f
	getIpconfig.Stdout = f

	helpers.WriteLog(nameOfContentFile, "----------------------------------------------------SYSTEMINFO------------------------------------------------\n", 0)
	err = getSystemInfo.Run()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(3)
	}

	helpers.WriteLog(nameOfContentFile, "----------------------------------------------------SYSTEMINFOFROMPOWERSHELL------------------------------------------------\n", 0)
	err = getPowershellSystemInfo.Run()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(4)
	}

	helpers.WriteLog(nameOfContentFile, "----------------------------------------------------TASKLIST------------------------------------------------\n", 0)
	err = getTasklist.Run()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(5)
	}

	helpers.WriteLog(nameOfContentFile, "----------------------------------------------------NETSTAT------------------------------------------------\n", 0)
	err = getNetstat.Run()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(6)
	}

	helpers.WriteLog(nameOfContentFile, "----------------------------------------------------PORTOCOLSTATISTICS------------------------------------------------\n", 0)
	err = getPortocolStatistics.Run()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(7)
	}

	helpers.WriteLog(nameOfContentFile, "----------------------------------------------------IPCONFIG------------------------------------------------\n", 0)
	err = getIpconfig.Run()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(8)
	}
	helpers.WriteLog(nameOfLogFile, "Ending Test: GetSysInfo", 2)
	os.Exit(0)

}
