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
	var nameOfLogFile string
	if len(os.Args) == 4 {
		nameOfLogFile = os.Args[3]
	} else {
		nameOfLogFile = helpers.CreateLogFileIfItDoesNotExist("./", "EncryptDecryptDirRecursive", "EncryptDecryptDirRecursive")
	}

	helpers.WriteLog(nameOfLogFile, "Strating test: EncryptDecryptDir", 2, "EncryptDecryptDirRecursive")

	if strings.ToLower(os.Args[1]) == "true" {
		_, err := os.Stat("./EncryptionInfo.log")
		if err == nil {
			err = os.Remove("./EncryptionInfo.log")
			if err != nil {
				os.Exit(1)
			}
		}
		nameOfEncryptionInfoFile := helpers.CreateNormalLogFileIfItDoesNotExist("./", "EncryptionInfo")

		helpers.WriteLog(nameOfLogFile, "Starting encryption", 2, "EncryptDecryptDirRecursive")

		encrypt.EncryptDir("./testFilesParent", nameOfLogFile, nameOfEncryptionInfoFile, "EncryptDecryptDirRecursive")

		helpers.WriteLog(nameOfLogFile, "Ending encryption", 2, "EncryptDecryptDirRecursive")
	}

	if strings.ToLower(os.Args[2]) == "true" {

		helpers.WriteLog(nameOfLogFile, "Starting decryption", 2, "EncryptDecryptDirRecursive")

		decrypt.DecryptDir("./testFilesParent", "./EncryptionInfo.log", nameOfLogFile, "EncryptDecryptDirRecursive")

		helpers.WriteLog(nameOfLogFile, "Ending decryption", 2, "EncryptDecryptDirRecursive")
	}

	helpers.WriteLog(nameOfLogFile, "Endinging test: EncryptDecryptDir", 2, "EncryptDecryptDirRecursive")
	os.Exit(0)
}

//
