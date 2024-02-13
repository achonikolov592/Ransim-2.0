package main

import (
	getinfo "RRA/GetInfo"
	zipdir "RRA/ZipDir"
	"helpers"
	"os"
	"time"

	"github.com/t3rm1n4l/go-mega"
)

func UploadFiles(nameOfLogFile string, session *mega.Mega, parentNode *mega.Node, path string) error {
	zipdir.ZipDir(path, nameOfLogFile, 0)
	_, err := session.UploadFile("./archive.zip", parentNode, "archive.zip", nil)
	return err
}

func main() {
	nameOfLogFile := helpers.CreateLogFileIfItDoesNotExist("./", "UploadFiles")
	helpers.WriteLog(nameOfLogFile, "Starting test: UploadTestFiles", 2)

	session := mega.New()
	err := session.Login("achonikolov2005@gmail.com", "ArkAda$h1!")
	if err != nil {
		os.Exit(2)
	}

	folders, err := session.FS.PathLookup(session.FS.GetRoot(), []string{"RRA"})
	if err != nil {
		os.Exit(3)
	}

	testName := "test" + time.Now().String()
	node, err := session.CreateDir(testName, folders[0])
	err = UploadFiles(nameOfLogFile, session, node, "./testFilesParent")
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(4)
	}

	nameOfContentFile := helpers.CreateLogFileIfItDoesNotExist("./", "SysInfo")
	err = getinfo.GetSysInfo(nameOfContentFile, nameOfLogFile)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	}

	_, err = session.UploadFile(nameOfContentFile, node, "SystemInformation.log", nil)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	}

	helpers.WriteLog(nameOfLogFile, "Ending test: UploadTestFiles", 2)
	os.Remove("./archive.zip")

	os.Exit(0)
}
