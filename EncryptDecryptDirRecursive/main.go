package main

import (
	"RRA/EncryptDecryptDirRecursive/decrypt"
	"RRA/EncryptDecryptDirRecursive/encrypt"
	"helpers"
	"os"
	"strings"
)

var nameOfLogFile string

func main() {
	nameOfLogFile = helpers.CreateLogFileIfItDoesNotExist("./", "EncryptDecryptDirRecursive")

	helpers.WriteLog(nameOfLogFile, "Strating test: EncryptDecryptDir", 2)

	if strings.ToLower(os.Args[1]) == "true" {
		_, err := os.Stat("./EncryptionInfo.log")
		if err == nil {
			err = os.Remove("./EncryptionInfo.log")
			if err != nil {
				os.Exit(1)
			}
		}
		nameOfEncryptionInfoFile := helpers.CreateLogFileIfItDoesNotExist("./", "EncryptionInfo")

		helpers.WriteLog(nameOfLogFile, "Starting encryption", 2)

		encrypt.EncryptDir("./testFilesParent", nameOfLogFile, nameOfEncryptionInfoFile)

		helpers.WriteLog(nameOfLogFile, "Ending encryption", 2)
	}

	if strings.ToLower(os.Args[2]) == "true" {

		helpers.WriteLog(nameOfLogFile, "Starting decryption", 2)

		decrypt.DecryptDir("./testFilesParent", "./EncryptionInfo.log", nameOfLogFile)

		helpers.WriteLog(nameOfLogFile, "Ending decryption", 2)
	}

	helpers.WriteLog(nameOfLogFile, "Endinging test: EncryptDecryptDir", 2)
	os.Exit(0)
}

//
