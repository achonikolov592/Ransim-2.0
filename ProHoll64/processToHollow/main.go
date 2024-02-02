package main

import (
	"syscall"
	"unsafe"
)

func main() {
	syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		uintptr(0),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Hello world from 64"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("Hello world from 64!"))),
		uintptr(0))
}
