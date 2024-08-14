package main

import (
	zipdir "RRA/ZipDir"
	"helpers"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	var nameOfLogFile string
	if len(os.Args) == 3 {
		nameOfLogFile = os.Args[2]
	} else {
		nameOfLogFile, _ = filepath.Abs(helpers.CreateLogFileIfItDoesNotExist("./", "ArchiveFiles", "ArchiveFiles"))
	}

	helpers.WriteLog(nameOfLogFile, "Starting test: ArchiveFiles", 2, "ArchiveFiles")
	_, err := os.Stat("./archive.zip")
	if err == nil {
		os.Remove("./archive.zip")
	}

	timeToDelay, _ := strconv.Atoi(os.Args[1])
	zipdir.ZipDir("./testFilesParent", nameOfLogFile, timeToDelay, "ArchiveFiles")
	helpers.WriteLog(nameOfLogFile, "Ending test: ArchiveFiles", 2, "ArchiveFiles")
}

//
