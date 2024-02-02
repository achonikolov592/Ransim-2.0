package main

import "C"
import (
	"os"
)

//export DLLMain
func DLLMain() {
	//windows.MessageBox(windows.HWND(0), syscall.StringToUTF16Ptr("Injected"), syscall.StringToUTF16Ptr("Injection works"), windows.MB_OK)
	os.Create("C:\\users\\achon\\onedrive\\desktop\\a.txt")
}

func main() {}

//a
