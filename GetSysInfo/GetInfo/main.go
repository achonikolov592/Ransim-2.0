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

		helpers.WriteLog(nameOfContentFile, "----------------------------------------------------UNAME------------------------------------------------\n", 0)
		out, err := getuname.Output()
		helpers.WriteLog(nameOfContentFile, string(out), 0)
		if err != nil {
			return err
		}
		helpers.WriteLog(nameOfContentFile, "----------------------------------------------------GROPUS------------------------------------------------\n", 0)
		out, err = getGroups.Output()
		helpers.WriteLog(nameOfContentFile, string(out), 0)
		if err != nil {
			return err
		}
		helpers.WriteLog(nameOfContentFile, "----------------------------------------------------PROCESSES------------------------------------------------\n", 0)
		out, err = getProcesses.Output()
		helpers.WriteLog(nameOfContentFile, string(out), 0)
		if err != nil {
			return err
		}
		helpers.WriteLog(nameOfContentFile, "----------------------------------------------------NETSTAT------------------------------------------------\n", 0)
		out, err = getNetstat.Output()
		helpers.WriteLog(nameOfContentFile, string(out), 0)
		if err != nil {
			return err
		}
		helpers.WriteLog(nameOfContentFile, "----------------------------------------------------PROTOCOLSSTATISTICS------------------------------------------------\n", 0)
		out, err = getPortocolStatistics.Output()
		helpers.WriteLog(nameOfContentFile, string(out), 0)
		if err != nil {
			return err
		}
		helpers.WriteLog(nameOfContentFile, "----------------------------------------------------IFCONFIG------------------------------------------------\n", 0)
		out, err = getIfconfig.Output()
		helpers.WriteLog(nameOfContentFile, string(out), 0)
		if err != nil {
			return err
		}
	}

	return nil
}
