package main

import (
	"fmt"
	"log"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/windows"
)

func main() {
	dPath, _ := filepath.Abs("..\\DLLInjection\\dll\\mydll.dll")
	pId := uintptr(2696)

	kernel32 := windows.NewLazyDLL("kernel32.dll")

	pHandle, err := windows.OpenProcess(windows.PROCESS_CREATE_THREAD|windows.PROCESS_VM_OPERATION|windows.PROCESS_VM_WRITE|windows.PROCESS_VM_READ|windows.PROCESS_QUERY_INFORMATION, false, uint32(pId))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Process opened")

	VirtualAllocEx := kernel32.NewProc("VirtualAllocEx")
	vAlloc, _, err := VirtualAllocEx.Call(uintptr(pHandle), 0, uintptr(len(dPath)+1), windows.MEM_RESERVE|windows.MEM_COMMIT, windows.PAGE_EXECUTE_READWRITE)
	fmt.Println("Memory allocated")

	bPtrDpath, err := windows.BytePtrFromString(dPath)
	if err != nil {
		log.Fatal(err)
	}

	Zero := uintptr(0)
	err = windows.WriteProcessMemory(pHandle, vAlloc, bPtrDpath, uintptr(len(dPath)+1), &Zero)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("DLL path written")

	LoadLibAddr, err := syscall.GetProcAddress(syscall.Handle(kernel32.Handle()), "LoadLibraryA")
	if err != nil {
		log.Fatal(err)
	}

	tHandle, _, _ := kernel32.NewProc("CreateRemoteThread").Call(uintptr(pHandle), 0, 0, LoadLibAddr, vAlloc, 0, 0)
	defer syscall.CloseHandle(syscall.Handle(tHandle))
	fmt.Println("DLL Injected")
	////

}
