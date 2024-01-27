package main

import (
	"RRA/EncryptDecryptDirRecursivePartially/decrypt"
	"RRA/EncryptDecryptDirRecursivePartially/encrypt"
	"helpers"
	"os"
	"strings"
)

var nameOfLogFile string

func main() {
	nameOfLogFile = helpers.CreateLogFileIfItDoesNotExist("./", "EncryptDecryptDirRecursivePartially")

	if strings.ToLower(os.Args[1]) == "true" {

		_, err := os.Stat("./EncryptionInfo.log")
		if err == nil {
			err = os.Remove("./EncryptionInfo.log")
			if err != nil {
				os.Exit(1)
			}
		}

		nameOfInfoFile := helpers.CreateLogFileIfItDoesNotExist("./", "EncryptionInfo")
		helpers.RemoveTestFilesIfExists("./")
		//helpers.CreateMultipleTestFiles("./", nameOfLogFile)

		helpers.WriteLog(nameOfLogFile, "Strating test: EncryptDirPartially", 2)

		encrypt.EncryptFilesInDir("./testFilesParent", nameOfLogFile, nameOfInfoFile)

		helpers.WriteLog(nameOfLogFile, "Ending test: EncryptDecryptDirRecursivePartially", 2)
	}
	if strings.ToLower(os.Args[2]) == "true" {
		helpers.WriteLog(nameOfLogFile, "Strating test: DecryptDirRecursivePartially", 2)

		decrypt.DecryptDir("./testFilesParent", nameOfLogFile, "./EncryptionInfo.log")

		helpers.WriteLog(nameOfLogFile, "Ending test: DecryptDirRecursivePartially", 2)
	}

	helpers.WriteLog(nameOfLogFile, "Endinging test: EncryptDecryptDirPartially", 2)
	os.Exit(0)
}

//
