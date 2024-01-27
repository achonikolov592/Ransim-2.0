package main

import (
	"RRA/SecureDeleteFiles/SecureDeleteFile"
	"helpers"
	"os"
	"path/filepath"
	"strconv"
)

func deleteFilesInDir(dir string, nameOfLogFile string, timeToDelay int) {
	var filesInDir []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if path != dir {
			filesInDir = append(filesInDir, path)
		}
		return nil
	})

	for _, file := range filesInDir {
		info, err := os.Stat(file)
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1)
			os.Exit(2)
		}

		if !(info.IsDir()) {
			SecureDeleteFile.SecureDelete(file, nameOfLogFile, timeToDelay)
		} else {
			deleteFilesInDir(file, nameOfLogFile, timeToDelay)
		}
	}

	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(3)
	}

	err = os.Remove(dir)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(4)
	}
}

func main() {
	nameOfLogFile := helpers.CreateLogFileIfItDoesNotExist("./", "SecureDeleteFiles")
	//helpers.CreateMultipleTestFiles("./", nameOfLogFile)

	helpers.WriteLog(nameOfLogFile, "Starting test: SecureDeleteFiles", 2)

	timeToDelay, _ := strconv.Atoi(os.Args[1])
	deleteFilesInDir("./testFilesParent", nameOfLogFile, timeToDelay)

	helpers.WriteLog(nameOfLogFile, "Ending test: SecureDeleteFiles", 2)

	os.Exit(0)
}

//
