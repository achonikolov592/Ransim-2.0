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

	if len(os.Args) == 4 {
		nameOfLogFile = os.Args[3]
	} else {
		nameOfLogFile = helpers.CreateLogFileIfItDoesNotExist("./", "EncryptDecryptDirRecursivePartially", "EncryptDecryptDirRecursivePartially")
	}

	if strings.ToLower(os.Args[1]) == "true" {

		_, err := os.Stat("./EncryptionInfo.log")
		if err == nil {
			err = os.Remove("./EncryptionInfo.log")
			if err != nil {
				os.Exit(1)
			}
		}

		nameOfInfoFile := helpers.CreateNormalLogFileIfItDoesNotExist("./", "EncryptionInfo")
		helpers.WriteLog(nameOfLogFile, "Strating test: EncryptDirPartially", 2, "EncryptDecryptDirRecursivePartially")

		encrypt.EncryptFilesInDir("./testFilesParent", nameOfLogFile, nameOfInfoFile, "EncryptDecryptDirRecursivePartially")

		helpers.WriteLog(nameOfLogFile, "Ending test: EncryptDecryptDirRecursivePartially", 2, "EncryptDecryptDirRecursivePartially")
	}
	if strings.ToLower(os.Args[2]) == "true" {
		helpers.WriteLog(nameOfLogFile, "Strating test: DecryptDirRecursivePartially", 2, "EncryptDecryptDirRecursivePartially")

		decrypt.DecryptDir("./testFilesParent", nameOfLogFile, "./EncryptionInfo.log", "EncryptDecryptDirRecursivePartially")

		helpers.WriteLog(nameOfLogFile, "Ending test: DecryptDirRecursivePartially", 2, "EncryptDecryptDirRecursivePartially")
	}

	helpers.WriteLog(nameOfLogFile, "Endinging test: EncryptDecryptDirPartially", 2, "EncryptDecryptDirRecursivePartially")
	os.Exit(0)
}

//
