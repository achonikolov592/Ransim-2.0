package main

import (
	"RRA/SecureDeleteFiles/SecureDeleteFile"
	"archive/zip"
	"helpers"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
)

func zipDir(dir string, nameOfLogFile string, timeToDelay int) {
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
			SecureDeleteFile.SecureDelete(file, nameOfLogFile, timeToDelay)
		}

		archivedFile.Close()
	}

	zipWriting.Close()
	os.RemoveAll(dir)
}
func main() {
	nameOfLogFile, _ := filepath.Abs(helpers.CreateLogFileIfItDoesNotExist("./", "Archive"))
	//helpers.CreateMultipleTestFiles("./", nameOfLogFile)

	helpers.WriteLog(nameOfLogFile, "Starting test: ArchiveFiles", 2)
	_, err := os.Stat("./archive.zip")
	if err == nil {
		os.Remove("./archive.zip")
	}

	timeToDelay, _ := strconv.Atoi(os.Args[1])
	zipDir("./testFilesParent", nameOfLogFile, timeToDelay)
	helpers.WriteLog(nameOfLogFile, "Ending test: ArchiveFiles", 2)
}

//
