package main

import (
	zipdir "RRA/ZipDir"
	"fmt"
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

	session := mega.New()
	err := session.Login("achonikolov2005@gmail.com", "ArkAda$h1!")
	if err != nil {
		fmt.Println(err)
	}
	folders, err := session.FS.PathLookup(session.FS.GetRoot(), []string{"RRA"})
	if err != nil {
		fmt.Println(err)
	}

	testName := "test" + time.Now().String()
	node, err := session.CreateDir(testName, folders[0])
	err = UploadFiles(nameOfLogFile, session, node, "./testFilesParent")
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	}

	os.Remove("./archive.zip")

}
