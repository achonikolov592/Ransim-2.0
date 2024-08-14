package main

import (
	"RRA/SecureDeleteFiles/SecureDeleteFile"
	"helpers"
	"os"
	"path/filepath"
	"strconv"
)

func deleteFilesInDir(dir string, nameOfLogFile string, timeToDelay int, nameOfTest string) {
	var filesInDir []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if path != dir {
			filesInDir = append(filesInDir, path)
		}
		return nil
	})

	for i := 0; i < len(filesInDir); i++ {
		info, err := os.Stat(filesInDir[i])
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1, nameOfTest)
			os.Exit(2)
		}

		if !(info.IsDir()) {
			SecureDeleteFile.SecureDelete(filesInDir[i], nameOfLogFile, timeToDelay, nameOfTest)
		}
	}

	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, nameOfTest)
		os.Exit(3)
	}

	err = os.RemoveAll(dir)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, nameOfTest)
		os.Exit(4)
	}
}

func main() {
	var nameOfLogFile string
	if len(os.Args) == 3 {
		nameOfLogFile = os.Args[2]
	} else {
		nameOfLogFile = helpers.CreateLogFileIfItDoesNotExist("./", "SecureDeleteFiles", "SecureDeleteFiles")
	}

	helpers.WriteLog(nameOfLogFile, "Starting test: SecureDeleteFiles", 2, "SecureDeleteFiles")

	timeToDelay, _ := strconv.Atoi(os.Args[1])
	deleteFilesInDir("./testFilesParent", nameOfLogFile, timeToDelay, "SecureDeleteFiles")

	helpers.WriteLog(nameOfLogFile, "Ending test: SecureDeleteFiles", 2, "SecureDeleteFiles")

	os.Exit(0)
}

//aaaaaa
