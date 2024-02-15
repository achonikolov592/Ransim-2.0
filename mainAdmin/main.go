package main

import (
	//obfuscate "RRA/Obfuscation"
	"encoding/json"
	"fmt"
	"helpers"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type customError struct {
	Message string
}

func (e *customError) Error() string {
	return e.Message
}

type location struct {
	name           string
	path           string
	nameOfFile     string
	expectedResult string
}

type Test struct {
	Name      string
	IsEnabled bool
}

type Setts struct {
	DirNumber                 int
	DirNumberFiles            int
	DocumentFiles             string
	PictureFiles              string
	TimeToDelayOnSecureDelete int
	ToEncrypt                 string
	ToDecrypt                 string
	MegaEmail                 string
	MegaPassword              string
	IPAddrForReverseShell     string
}

type Config struct {
	Tests    []Test
	Settings Setts
}

var testLocation = []location{location{"ArchiveFiles", "../ArchiveFiles/", "ArchiveFiles.exe", "nil"},
	location{"Eicar", "../Eicar/", "Eic.exe", "exit status 1"},
	location{"EncryptDecryptDirRecursive", "../EncryptDecryptDirRecursive/", "EncryptDecryptDirRecursive.exe", "nil"},
	location{"EncryptDecryptDirRecursivePartially", "../EncryptDecryptDirRecursivePartially/", "EncryptDecryptDirRecursivePartially.exe", "nil"},
	location{"GetSystemInformation", "../GetSysInfo/", "GetSysInfo.exe", "nil"},
	location{"SecureDeleteFiles", "../SecureDeleteFiles/", "SecureDeleteFiles.exe", "nil"},
	location{"StartupFolderNewFile", "../StartupFolderNewFile/", "Startup.exe", "exit status 1"},
	location{"ServiceCreation", "../ServiceCreation/", "ServiceCreation.exe", "exit status 1"},
	location{"PrivilegeEscalation", "../PrivEsc/", "AccsTok.exe", "exit status 1"},
	location{"RansomwareNoteDeploy", "../RanNote/", "RanNote.exe", "nil"},
	location{"ExfiltrationOfFiles", "../UploadFiles/", "Exfilt.exe", "nil"},
	location{"ProcessHollowing32", "../ProHoll32/", "ProHoll32.exe", "exit status 1"},
	location{"ProcessHollowing64", "../ProHoll64/", "ProHoll64.exe", "exit status 1"},
	location{"RegistryKeysTest", "../RegKeys/", "RegKeys.exe", "exit status 1"},
	location{"DLLSideLoading", "../DLLSideLoading/", "DllLoad.exe", "exit status 1"},
	location{"ReverseShell", "../RevSh/", "RevSh.exe", "exit status 1"},
	location{"OSCredentialDump", "../OSCredDump/", "CredDump.exe", "exit status 1"}}

var locationForTestFolder = []string{"../ArchiveFiles", "../EncryptDecryptDirRecursive", "../EncryptDecryptDirRecursivePartially", "../SecureDeleteFiles", "../ServiceCreation", "../UploadFiles"}

func getTestConfig() (Config, error) {
	var config Config
	jsonFile, err := os.Open("./tests.json")
	if err != nil {
		return config, err
	}

	value, _ := io.ReadAll(jsonFile)
	json.Unmarshal(value, &config)

	return config, nil
}

func findLocationByName(name string) (location, error) {

	for i := 0; i < len(testLocation); i++ {
		if name == testLocation[i].name {
			return testLocation[i], nil
		}
	}

	var loc location
	return loc, &customError{Message: "Invalid name"}
}

func prepareTestFiles(nameOfLogFile string, settings Setts) {

	DocFilesFormat := strings.Fields(settings.DocumentFiles)
	PicFilesFormat := strings.Fields(settings.PictureFiles)

	var files []*os.File
	err := filepath.Walk("../TestFilesAll", func(path string, info fs.FileInfo, err error) error {
		if !(info.IsDir()) {
			i := 0
			for i = 0; i < len(DocFilesFormat); i++ {
				if path[len(path)-len(DocFilesFormat[i]):] == DocFilesFormat[i] {
					originalFile, err := os.OpenFile(path, os.O_RDONLY, 0666)
					if err != nil {
						return err
					}
					files = append(files, originalFile)
				}
			}
			if i == len(DocFilesFormat) {
				for i = 0; i < len(PicFilesFormat); i++ {
					if path[len(path)-len(PicFilesFormat[i]):] == PicFilesFormat[i] {
						originalFile, err := os.OpenFile(path, os.O_RDONLY, 0666)
						if err != nil {
							return err
						}
						files = append(files, originalFile)
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(1)
	}

	for i := 0; i < len(locationForTestFolder); i++ {
		helpers.CreateMultipleTestFiles(locationForTestFolder[i], nameOfLogFile, files, settings.DirNumber, settings.DirNumberFiles)
	}

}

func getWhichTestToExecuteAndSettings(nameOfLogFile string) ([]location, Setts) {
	config, err := getTestConfig()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(1)
	}

	var whichTests []location
	for i := 0; i < len(config.Tests); i++ {
		if config.Tests[i].IsEnabled {
			if config.Tests[i].Name[:3] == "Win" || config.Tests[i].Name[:3] == "Lin" {
				if strings.ToLower(config.Tests[i].Name[:3]) == runtime.GOOS[:3] {
					loc, err := findLocationByName(config.Tests[i].Name[3:])
					if err != nil {
						fmt.Println(err.Error())
						os.Exit(1)
					}
					whichTests = append(whichTests, loc)
				}

			} else {
				loc, err := findLocationByName(config.Tests[i].Name[3:])
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				whichTests = append(whichTests, loc)
			}

		}
	}

	return whichTests, config.Settings
}

func isTestCorrect(name string, correctTests []string) bool {
	for i := 0; i < len(correctTests); i++ {
		if name == correctTests[i] {
			return true
		}
	}

	return false
}

func getResults(nameOfLogFile string, whichTest []location, correctTests []string) {
	os.Remove("./results.html")
	f, err := os.Create("./results.html")
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	}

	var content = "<html lang='en'>" +
		"<style>" +
		"table, th, td{border:1px solid black}" +
		"</style>" +
		"<head>" +
		"<meta charset='UTF-8'>" +
		"<meta name='viewport' content='width=device-width, initial-scale=1.0'>" +
		"<title>Results</title>" +
		"</head>" +
		"<body>" +
		"<h1> Table of results <h1>" +
		"<table>" +
		"<tr>" + "<th>Name Of Test</th>" + "<th>Is correct</th>" + "</tr>"
	for i := 0; i < len(whichTest); i++ {
		if isTestCorrect(whichTest[i].name, correctTests) {
			content += "<tr>" + "<td style='background-color: lightgreen'>" + whichTest[i].name + "</td>" + "<td style='background-color: lightgreen' >Yes</td>" + "<tr>"
		} else {
			content += "<tr>" + "<td style='background-color: LightCoral'>" + whichTest[i].name + "</td>" + "<td style='background-color: LightCoral' >No</td>" + "<tr>"
		}
	}

	content += "</table>" + "</body>" + "</html>"
	f.WriteString(content)
}

func main() {
	nameOfLogFile := helpers.CreateLogFileIfItDoesNotExist("./", "main")

	whichTests, settings := getWhichTestToExecuteAndSettings("tests.json")
	prepareTestFiles(nameOfLogFile, settings)
	//obfuscate.PrepareEveryTestObfuscated(nameOfLogFile)

	var correctTests, incorrectTests []string
	for i := 0; i < len(whichTests); i++ {
		fmt.Println("Starting test: " + whichTests[i].name)
		if whichTests[i].name == "EncryptDecryptDirRecursive" || whichTests[i].name == "EncryptDecryptDirRecursivePartially" {
			cmd := exec.Command(whichTests[i].path+whichTests[i].nameOfFile, settings.ToEncrypt, settings.ToDecrypt)
			cmd.Dir = whichTests[i].path
			err := cmd.Run()

			if err == nil {
				correctTests = append(correctTests, whichTests[i].name)
			} else {
				helpers.WriteLog(nameOfLogFile, err.Error()+" from "+whichTests[i].name, 1)
				incorrectTests = append(incorrectTests, whichTests[i].name)
			}
		} else if whichTests[i].name == "ServiceCreation" {
			cmd := exec.Command(whichTests[i].path+whichTests[i].nameOfFile, "install")
			cmd.Dir = whichTests[i].path
			err := cmd.Run()

			if err == nil {
				if whichTests[i].expectedResult == "nil" {
					correctTests = append(correctTests, whichTests[i].name)
				} else {
					incorrectTests = append(incorrectTests, whichTests[i].name)
				}
				cmd := exec.Command(whichTests[i].path+whichTests[i].nameOfFile, "uninstall")
				cmd.Dir = whichTests[i].path
				cmd.Run()
			} else {
				if err.Error() == whichTests[i].expectedResult {
					correctTests = append(correctTests, whichTests[i].name)
				} else {
					helpers.WriteLog(nameOfLogFile, err.Error()+" from "+whichTests[i].name, 1)
					incorrectTests = append(incorrectTests, whichTests[i].name)
				}
			}
		} else if whichTests[i].name == "SecureDeleteFiles" || whichTests[i].name == "ArchiveFiles" {
			cmd := exec.Command(whichTests[i].path+whichTests[i].nameOfFile, strconv.Itoa(settings.TimeToDelayOnSecureDelete))
			cmd.Dir = whichTests[i].path
			err := cmd.Run()

			if err == nil {
				if whichTests[i].expectedResult == "nil" {
					correctTests = append(correctTests, whichTests[i].name)
				} else {
					incorrectTests = append(incorrectTests, whichTests[i].name)
				}
			} else {
				if err.Error() == whichTests[i].expectedResult {
					correctTests = append(correctTests, whichTests[i].name)
				} else {
					helpers.WriteLog(nameOfLogFile, err.Error()+" from "+whichTests[i].name, 1)
					incorrectTests = append(incorrectTests, whichTests[i].name)
				}
			}
		} else if whichTests[i].name == "ProcessHollowing32" {
			cmd := exec.Command(whichTests[i].path+whichTests[i].nameOfFile, "C:\\windows\\syswow64\\notepad.exe")
			cmd.Dir = whichTests[i].path
			err := cmd.Run()
			if err == nil {
				incorrectTests = append(incorrectTests, whichTests[i].name)
			} else {
				if err.Error() == whichTests[i].expectedResult {
					correctTests = append(correctTests, whichTests[i].name)
				} else {
					helpers.WriteLog(nameOfLogFile, err.Error()+" from "+whichTests[i].name, 1)
					incorrectTests = append(incorrectTests, whichTests[i].name)
				}
			}
		} else if whichTests[i].name == "ProcessHollowing64" {
			cmd := exec.Command(whichTests[i].path+whichTests[i].nameOfFile, "C:\\windows\\System32\\notepad.exe")
			cmd.Dir = whichTests[i].path
			err := cmd.Run()
			if err == nil {
				incorrectTests = append(incorrectTests, whichTests[i].name)
			} else {
				if err.Error() == whichTests[i].expectedResult {
					correctTests = append(correctTests, whichTests[i].name)
				} else {
					helpers.WriteLog(nameOfLogFile, err.Error()+" from "+whichTests[i].name, 1)
					incorrectTests = append(incorrectTests, whichTests[i].name)
				}
			}
		} else if whichTests[i].name == "ReverseShell" {
			cmd := exec.Command(whichTests[i].path+whichTests[i].nameOfFile, settings.IPAddrForReverseShell)
			cmd.Dir = whichTests[i].path
			err := cmd.Run()

			if err == nil {
				helpers.WriteLog(nameOfLogFile, err.Error()+" from "+whichTests[i].name, 1)
				incorrectTests = append(incorrectTests, whichTests[i].name)
			} else {
				if err.Error() == whichTests[i].expectedResult {
					correctTests = append(correctTests, whichTests[i].name)
				} else {
					incorrectTests = append(incorrectTests, whichTests[i].name)
				}
			}
		} else if whichTests[i].name == "ExfiltrationOfFiles" {
			cmd := exec.Command(whichTests[i].path+whichTests[i].nameOfFile, settings.MegaEmail, settings.MegaPassword)
			cmd.Dir = whichTests[i].path
			err := cmd.Run()

			if err == nil {
				correctTests = append(correctTests, whichTests[i].name)
			} else {
				helpers.WriteLog(nameOfLogFile, err.Error()+" from "+whichTests[i].name, 1)
				incorrectTests = append(incorrectTests, whichTests[i].name)
			}
		} else {
			cmd := exec.Command(whichTests[i].path + whichTests[i].nameOfFile)
			cmd.Dir = whichTests[i].path
			err := cmd.Run()

			if err == nil {
				if whichTests[i].expectedResult == "nil" {
					correctTests = append(correctTests, whichTests[i].name)
				} else {
					incorrectTests = append(incorrectTests, whichTests[i].name)
				}
			} else {
				if err.Error() == whichTests[i].expectedResult {
					correctTests = append(correctTests, whichTests[i].name)
				} else {
					helpers.WriteLog(nameOfLogFile, err.Error()+" from "+whichTests[i].name, 1)
					incorrectTests = append(incorrectTests, whichTests[i].name)
				}
			}
		}
	}

	for i := 0; i < len(whichTests); i++ {
		_, err := os.Open(whichTests[i].path + whichTests[i].nameOfFile)
		if err != nil {
			if whichTests[i].expectedResult == "exit status 1" {
				for j := 0; j < len(incorrectTests); j++ {
					if incorrectTests[j] == whichTests[i].name {
						incorrectTests = append(incorrectTests[:j], incorrectTests[j+1:]...)
						correctTests = append(correctTests, whichTests[i].name)
						break
					}
				}
			} else {
				for j := 0; j < len(correctTests); j++ {
					if correctTests[j] == whichTests[i].name {
						correctTests = append(correctTests[:j], correctTests[j+1:]...)
						incorrectTests = append(incorrectTests, whichTests[i].name)
						break
					}
				}
			}
		}
	}

	getResults(nameOfLogFile, whichTests, correctTests)
	fmt.Println("You can check the results!")

	for i := 0; i < len(locationForTestFolder); i++ {
		err := helpers.RemoveTestFilesIfExists(locationForTestFolder[i])
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		}
	}

}
