package main

import (
	"helpers"
	"os"
	"os/exec"

	"golang.org/x/sys/windows/registry"
)

func compileDLLs(nameOfLogFile string) {
	cmd := exec.Command("go", "build", "-o", "../DLLSideLoading/kernel32.dll", "-buildmode=c-shared", "dll.go")
	cmd1 := exec.Command("go", "build", "-o", "../DLLSideLoading/user32.dll", "-buildmode=c-shared", "dll.go")
	cmd2 := exec.Command("go", "build", "-o", "../DLLSideLoading/advapi32.dll", "-buildmode=c-shared", "dll.go")
	cmd3 := exec.Command("go", "build", "-o", "../DLLSideLoading/ole32.dll", "-buildmode=c-shared", "dll.go")
	cmd4 := exec.Command("go", "build", "-o", "../DLLSideLoading/shell32.dll", "-buildmode=c-shared", "dll.go")

	cmd.Dir = "../DLLs"
	cmd1.Dir = "../DLLs"
	cmd2.Dir = "../DLLs"
	cmd3.Dir = "../DLLs"
	cmd4.Dir = "../DLLs"

	err := cmd.Run()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	}
	err = cmd1.Run()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	}
	err = cmd2.Run()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	}
	err = cmd3.Run()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	}
	err = cmd4.Run()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	}
}

func main() {
	nameOfLogFile := helpers.CreateLogFileIfItDoesNotExist("./", "DLLSideLoad")
	numberOfErrors := 0
	numberOfNonCompiledDLLs := 0
	_, err := registry.OpenKey(registry.LOCAL_MACHINE, "System\\CurrentControlSet\\Control\\Session Manager", registry.WRITE)
	if err != nil {
		numberOfErrors += 1
	}

	compileDLLs(nameOfLogFile)
	_, err = os.OpenFile("./kernel32.dll", os.O_RDONLY, 0666)
	if err != nil {
		numberOfNonCompiledDLLs += 1
	}
	_, err = os.OpenFile("./advapi32.dll", os.O_RDONLY, 0666)
	if err != nil {
		numberOfNonCompiledDLLs += 1
	}
	_, err = os.OpenFile("./user32.dll", os.O_RDONLY, 0666)
	if err != nil {
		numberOfNonCompiledDLLs += 1
	}
	_, err = os.OpenFile("./ole32.dll", os.O_RDONLY, 0666)
	if err != nil {
		numberOfNonCompiledDLLs += 1
	}
	_, err = os.OpenFile("./shell32.dll", os.O_RDONLY, 0666)
	if err != nil {
		numberOfNonCompiledDLLs += 1
	}

	if numberOfNonCompiledDLLs == 5 {
		numberOfErrors += 1
	}

	if numberOfErrors == 2 {
		os.Exit(1)
	}

	os.Exit(0)
}
