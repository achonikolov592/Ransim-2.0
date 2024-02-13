package main

import "golang.org/x/sys/windows"

func main() {
	windows.MessageBox(0, windows.StringToUTF16Ptr("You are Encrypted!"), windows.StringToUTF16Ptr("Suprise"), 0)
}

//a
