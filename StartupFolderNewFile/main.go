package main

import (
	"helpers"
	"io"
	"os"
	"os/exec"
	"strings"
)

func main() {
	var name string
	if len(os.Args) == 2 {
		name = os.Args[1]
	} else {
		name = helpers.CreateLogFileIfItDoesNotExist("./", "StartupFolderNewFile", "StartupFolderNewFile")
	}

	helpers.WriteLog(name, "Strating test : StartupFolderNewFile", 2, "StartupFolderNewFile")

	src, err := os.Open("./ShowBox/ShowBox.exe")
	if err != nil {
		helpers.WriteLog(name, err.Error(), 1, "StartupFolderNewFile")
		os.Exit(3)
	}

	cmd := exec.Command("reg", "query", "HKEY_LOCAL_MACHINE\\SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Explorer\\Shell Folders", "/v", "Common Startup")
	result, err := cmd.Output()
	if err != nil {
		helpers.WriteLog(name, err.Error(), 1, "StartupFolderNewFile")
		os.Exit(4)
	}
	words := strings.Fields(string(result))

	dest, err := os.Create(words[len(words)-2] + " " + words[len(words)-1] + "/" + "a.exe")
	if err != nil {
		helpers.WriteLog(name, err.Error(), 1, "StartupFolderNewFile")
		os.Exit(1)
	} else {
		_, err = io.Copy(dest, src)
		if err != nil {
			helpers.WriteLog(name, err.Error(), 1, "StartupFolderNewFile")
			os.Exit(1)
		}

		_, err = os.Open(words[len(words)-2] + " " + words[len(words)-1] + "/" + "a.exe")
		if err != nil {
			helpers.WriteLog(name, err.Error(), 1, "StartupFolderNewFile")
			os.Exit(1)
		}

		err = os.Remove(words[len(words)-2] + " " + words[len(words)-1] + "/" + "a.exe")
		if err != nil {
			helpers.WriteLog(name, err.Error(), 1, "StartupFolderNewFile")
			os.Exit(7)
		}
	}

	helpers.WriteLog(name, "Ending test: StartupFolderNewFile", 2, "StartupFolderNewFile")
	os.Exit(0)
}

//aaa
