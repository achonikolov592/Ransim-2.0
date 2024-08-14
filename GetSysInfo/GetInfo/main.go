package getinfo

import (
	"helpers"
	"os"
	"os/exec"
	"runtime"
)

func GetSysInfo(nameOfContentFile string) error {
	if runtime.GOOS == "windows" {
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

		helpers.WriteNormalLog(nameOfContentFile, "----------------------------------------------------SYSTEMINFO------------------------------------------------\n")
		err = getSystemInfo.Run()
		if err != nil {
			return err
		}

		helpers.WriteNormalLog(nameOfContentFile, "----------------------------------------------------SYSTEMINFOFROMPOWERSHELL------------------------------------------------\n")
		err = getPowershellSystemInfo.Run()
		if err != nil {
			return err
		}

		helpers.WriteNormalLog(nameOfContentFile, "----------------------------------------------------TASKLIST------------------------------------------------\n")
		err = getTasklist.Run()
		if err != nil {
			return err
		}

		helpers.WriteNormalLog(nameOfContentFile, "----------------------------------------------------NETSTAT------------------------------------------------\n")
		err = getNetstat.Run()
		if err != nil {
			return err
		}

		helpers.WriteNormalLog(nameOfContentFile, "----------------------------------------------------PORTOCOLSTATISTICS------------------------------------------------\n")
		err = getPortocolStatistics.Run()
		if err != nil {
			return err
		}

		helpers.WriteNormalLog(nameOfContentFile, "----------------------------------------------------IPCONFIG------------------------------------------------\n")
		err = getIpconfig.Run()
		if err != nil {
			return err
		}

		helpers.WriteNormalLog(nameOfContentFile, "----------------------------------------------------NetBIOS------------------------------------------------\n")
		err = getNetBIOS.Run()
		if err != nil {
			return err
		}

		return nil
	} else {
		getuname := exec.Command("bash", "-c", "uname -a")
		getGroups := exec.Command("bash", "-c", "groups")
		getProcesses := exec.Command("bash", "-c", "ps -aux")
		getNetstat := exec.Command("bash", "-c", "netstat -v -a")
		getPortocolStatistics := exec.Command("bash", "-c", "netstat -s")
		getIfconfig := exec.Command("bash", "-c", "ifconfig -a")

		/*f, err := os.OpenFile(nameOfContentFile, os.O_APPEND, 0666)
		if err != nil {
			return err
		}*/

		helpers.WriteNormalLog(nameOfContentFile, "----------------------------------------------------UNAME------------------------------------------------\n")
		out, err := getuname.Output()
		helpers.WriteNormalLog(nameOfContentFile, string(out))
		if err != nil {
			return err
		}
		helpers.WriteNormalLog(nameOfContentFile, "----------------------------------------------------GROPUS------------------------------------------------\n")
		out, err = getGroups.Output()
		helpers.WriteNormalLog(nameOfContentFile, string(out))
		if err != nil {
			return err
		}
		helpers.WriteNormalLog(nameOfContentFile, "----------------------------------------------------PROCESSES------------------------------------------------\n")
		out, err = getProcesses.Output()
		helpers.WriteNormalLog(nameOfContentFile, string(out))
		if err != nil {
			return err
		}
		helpers.WriteNormalLog(nameOfContentFile, "----------------------------------------------------NETSTAT------------------------------------------------\n")
		out, err = getNetstat.Output()
		helpers.WriteNormalLog(nameOfContentFile, string(out))
		if err != nil {
			return err
		}
		helpers.WriteNormalLog(nameOfContentFile, "----------------------------------------------------PROTOCOLSSTATISTICS------------------------------------------------\n")
		out, err = getPortocolStatistics.Output()
		helpers.WriteNormalLog(nameOfContentFile, string(out))
		if err != nil {
			return err
		}
		helpers.WriteNormalLog(nameOfContentFile, "----------------------------------------------------IFCONFIG------------------------------------------------\n")
		out, err = getIfconfig.Output()
		helpers.WriteNormalLog(nameOfContentFile, string(out))
		if err != nil {
			return err
		}
	}

	return nil
}
