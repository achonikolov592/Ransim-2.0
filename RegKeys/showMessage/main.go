package main

import (
	"syscall"
	"unsafe"
)

func main() {
	syscall.NewLazyDLL("user32.dll").NewProc("MessageBoxW").Call(
		uintptr(0),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("You have been encrypted! Be more careful next time!"))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("You have been encrypted!"))),
		uintptr(0))
}
