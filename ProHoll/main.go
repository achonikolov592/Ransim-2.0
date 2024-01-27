package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	proHollExe, _ := filepath.Abs("./go_libpeconv")
	targetPath := os.Args[1]
	cmd := exec.Command(proHollExe, "./test.exe", targetPath)
	err := cmd.Run()
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
