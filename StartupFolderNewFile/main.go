package main

import (
	"helpers"
	"io"
	"os"
	"os/exec"
	"strings"
)

func main() {
	name := helpers.CreateLogFileIfItDoesNotExist("./", "startup")
	helpers.CreateLogFileIfItDoesNotExist("./", "EncryptionInfo")
	helpers.WriteLog(name, "Strating test : StartupFolderNewFile", 2)

	compileFile := exec.Command("go", "build", ".")
	compileFile.Dir = "./encr"

	err := compileFile.Run()
	if err != nil {
		helpers.WriteLog(name, err.Error(), 1)
		os.Exit(2)
	}

	src, err := os.Open("./encr/enc.exe")
	if err != nil {
		helpers.WriteLog(name, err.Error(), 1)
		os.Exit(3)
	}

	cmd := exec.Command("reg", "query", "HKEY_LOCAL_MACHINE\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\Shell Folders", "/v", "Common Startup")
	result, err := cmd.Output()
	if err != nil {
		helpers.WriteLog(name, err.Error(), 1)
		os.Exit(4)
	}
	words := strings.Fields(string(result))

	dest, err := os.Create(words[len(words)-2] + " " + words[len(words)-1] + "/" + "a.exe")
	if err != nil {
		helpers.WriteLog(name, err.Error(), 1)
		os.Exit(1)
	}
	_, err = io.Copy(dest, src)
	if err != nil {
		helpers.WriteLog(name, err.Error(), 1)
		os.Exit(5)
	}

	_, err = os.Open(words[len(words)-2] + " " + words[len(words)-1] + "/" + "a.exe")
	if err != nil {
		helpers.WriteLog(name, err.Error(), 1)
		os.Exit(1)
	}

	helpers.WriteLog(name, "Ending test: StartupFolderNewFile", 2)
	os.Exit(0)
}

//aaa
