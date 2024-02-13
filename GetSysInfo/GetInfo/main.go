package getinfo

import (
	"helpers"
	"os"
	"os/exec"
)

func GetSysInfo(nameOfContentFile, nameOfLogFile string) error {
	getSystemInfo := exec.Command("systeminfo")
	getPowershellSystemInfo := exec.Command("powershell.exe", "Get-ComputerInfo")
	getTasklist := exec.Command("tasklist")
	getNetstat := exec.Command("netstat", "-a")
	getPortocolStatistics := exec.Command("netstat", "-s")
	getIpconfig := exec.Command("ipconfig", "/all")
	getNetBIOS := exec.Command("nbtstat", "-S")

	f, err := os.OpenFile(nameOfContentFile, os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	getSystemInfo.Stdout = f
	getPowershellSystemInfo.Stdout = f
	getTasklist.Stdout = f
	getNetstat.Stdout = f
	getPortocolStatistics.Stdout = f
	getIpconfig.Stdout = f
	getNetBIOS.Stdout = f

	helpers.WriteLog(nameOfContentFile, "----------------------------------------------------SYSTEMINFO------------------------------------------------\n", 0)
	err = getSystemInfo.Run()
	if err != nil {
		return err
	}

	helpers.WriteLog(nameOfContentFile, "----------------------------------------------------SYSTEMINFOFROMPOWERSHELL------------------------------------------------\n", 0)
	err = getPowershellSystemInfo.Run()
	if err != nil {
		return err
	}

	helpers.WriteLog(nameOfContentFile, "----------------------------------------------------TASKLIST------------------------------------------------\n", 0)
	err = getTasklist.Run()
	if err != nil {
		return err
	}

	helpers.WriteLog(nameOfContentFile, "----------------------------------------------------NETSTAT------------------------------------------------\n", 0)
	err = getNetstat.Run()
	if err != nil {
		return err
	}

	helpers.WriteLog(nameOfContentFile, "----------------------------------------------------PORTOCOLSTATISTICS------------------------------------------------\n", 0)
	err = getPortocolStatistics.Run()
	if err != nil {
		return err
	}

	helpers.WriteLog(nameOfContentFile, "----------------------------------------------------IPCONFIG------------------------------------------------\n", 0)
	err = getIpconfig.Run()
	if err != nil {
		return err
	}

	helpers.WriteLog(nameOfContentFile, "----------------------------------------------------NetBIOS------------------------------------------------\n", 0)
	err = getNetBIOS.Run()
	if err != nil {
		return err
	}

	return nil
}
