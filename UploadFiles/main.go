package main

import (
	"archive/zip"
	"fmt"
	"helpers"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/t3rm1n4l/go-mega"
)

func zipDir(dir string, nameOfLogFile string) {
	var files []string

	errFileWalk := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	if errFileWalk != nil {
		helpers.WriteLog(nameOfLogFile, errFileWalk.Error(), 1)
		os.Exit(2)
	}

	archive, err := os.Create("./archive.zip")
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(3)
	}
	defer archive.Close()
	zipWriting := zip.NewWriter(archive)

	for _, file := range files {
		archivedFile, errOfOpeningFile := os.Open(file)
		if errOfOpeningFile != nil {
			helpers.WriteLog(nameOfLogFile, errOfOpeningFile.Error(), 1)
			os.Exit(4)
		}

		fileInfo, _ := os.Stat(file)

		if !fileInfo.IsDir() {
			w1, err := zipWriting.Create(file)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
				os.Exit(5)
			}
			if _, err := io.Copy(w1, archivedFile); err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
				os.Exit(6)
			}
			//fileAbs, err := filepath.Abs(file)
			//fmt.Println(nameOfLogFile)
			//SecureDeleteFile.SecureDelete(fileAbs, nameOfLogFile)
		}

		archivedFile.Close()
	}

	zipWriting.Close()
}

func UploadFiles(nameOfLogFile string, session *mega.Mega, parentNode *mega.Node, path string) error {
	zipDir(path, nameOfLogFile)
	_, err := session.UploadFile("./archive.zip", parentNode, "archive.zip", nil)
	return err
}

func main() {
	nameOfLogFile := helpers.CreateLogFileIfItDoesNotExist("./", "UploadFiles")
	//helpers.CreateMultipleTestFiles("./", nameOfLogFile)

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

}
