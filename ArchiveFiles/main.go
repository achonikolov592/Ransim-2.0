package main

import (
	zipdir "RRA/ZipDir"
	"helpers"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	nameOfLogFile, _ := filepath.Abs(helpers.CreateLogFileIfItDoesNotExist("./", "Archive"))
	//helpers.CreateMultipleTestFiles("./", nameOfLogFile)

	helpers.WriteLog(nameOfLogFile, "Starting test: ArchiveFiles", 2)
	_, err := os.Stat("./archive.zip")
	if err == nil {
		os.Remove("./archive.zip")
	}

	timeToDelay, _ := strconv.Atoi(os.Args[1])
	zipdir.ZipDir("./testFilesParent", nameOfLogFile, timeToDelay)
	helpers.WriteLog(nameOfLogFile, "Ending test: ArchiveFiles", 2)
}

//
