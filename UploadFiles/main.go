package main

import (
	getinfo "RRA/GetInfo"
	zipdir "RRA/ZipDir"
	"helpers"
	"os"
	"time"

	mega "RRA/gomega"
)

func UploadFiles(nameOfLogFile string, session *mega.Mega, parentNode *mega.Node, path string, nameOfTest string) error {
	zipdir.ZipDir(path, nameOfLogFile, 0, nameOfTest)
	_, err := session.UploadFile("./archive.zip", parentNode, "archive.zip", nil)
	return err
}

func main() {
	var nameOfLogFile string
	if len(os.Args) == 4 {
		nameOfLogFile = os.Args[3]
	} else {
		nameOfLogFile = helpers.CreateLogFileIfItDoesNotExist("./", "UploadFiles", "UploadFiles")
	}

	helpers.WriteLog(nameOfLogFile, "Starting test: UploadTestFiles", 2, "UploadFiles")

	session := mega.New()
	if len(os.Args[1]) == 1 && len(os.Args[2]) == 1 {
		err := session.Login("achonikolov2005@gmail.com", "ArkAda$h1!")
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1, "UploadFiles")
			os.Exit(2)
		}
	} else {
		err := session.Login(os.Args[1], os.Args[2])
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1, "UploadFiles")
			os.Exit(2)
		}
	}

	folders, err := session.FS.PathLookup(session.FS.GetRoot(), []string{"RRA"})
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "UploadFiles")
		os.Exit(3)
	}

	testFilesName := "testFilesName" + time.Now().String()
	node, err := session.CreateDir(testFilesName, folders[0])
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "UploadFiles")
	}
	err = UploadFiles(nameOfLogFile, session, node, "./testFilesParent", "UploadFiles")
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "UploadFiles")
		os.Exit(4)
	}

	nameOfContentFile := helpers.CreateNormalLogFileIfItDoesNotExist("./", "SysInfo")
	err = getinfo.GetSysInfo(nameOfContentFile)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "UploadFiles")
		os.Exit(5)
	}

	_, err = session.UploadFile(nameOfContentFile, node, "SystemInformation.log", nil)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "UploadFiles")
		os.Exit(6)
	}

	helpers.WriteLog(nameOfLogFile, "Ending test: UploadTestFiles", 2, "UploadFiles")
	os.Remove("./archive.zip")

	os.Exit(0)
}
