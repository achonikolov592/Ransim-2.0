package main

import (
	obfuscate "RRA/Obfuscation"
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

type tests struct {
	Name      string
	IsEnabled bool
}

type DirSettings struct {
	DirNumber                 int
	DirNumberFiles            int
	DocumentFiles             string
	PictureFiles              string
	TimeToDelayOnSecureDelete int
	ToEncrypt                 string
	ToDecrypt                 string
}

var testLocation = []location{location{"ArchiveFiles", "../ArchiveFiles/", "ArchiveFiles.exe", "nil"},
	location{"Eicar", "../Eicar/", "Eic.exe", "exit status 1"},
	location{"EncryptDecryptDirRecursive", "../EncryptDecryptDirRecursive/", "EncryptDecryptDirRecursive.exe", "nil"},
	location{"EncryptDecryptDirRecursivePartially", "../EncryptDecryptDirRecursivePartially/", "EncryptDecryptDirRecursivePartially.exe", "nil"},
	location{"GetSystemInformation", "../GetSysInfo/", "GetSysInfo.exe", "nil"},
	location{"SecureDeleteFiles", "../SecureDeleteFiles/", "SecureDeleteFiles.exe", "nil"},
	location{"StartupFolderNewFile", "../StartupFolderNewFile/", "Startup.exe", "exit status 1"},
	location{"DownloadMaliciousPayload", "../MaliciousPayloadDownload/", "MaliciousPayloadDownload.exe", "exit status 1"},
	location{"ServiceCreation", "../ServiceCreation/", "ServiceCreation.exe", "exit status 1"},
	location{"PrivilegeEscalation", "../PrivEsc/", "AccsTok.exe", "exit status 1"},
	location{"RansomwareNoteDeploy", "../RanNote/", "RanNote.exe", "nil"},
	location{"ExfiltrationOfFIles", "../UploadFiles/", "Exfilt.exe", "nil"},
	location{"ProcessHollowing", "../ProHoll/", "ProHoll.exe", "exit status 1"}}

var locationForTestFolder = []string{"../ArchiveFiles", "../EncryptDecryptDirRecursive", "../EncryptDecryptDirRecursivePartially", "../SecureDeleteFiles", "../StartupFolderNewFile", "../UploadFiles"}

func getTestSettings() ([]tests, error) {
	var testSetts []tests
	jsonFile, err := os.Open("./tests.json")
	if err != nil {
		return testSetts, err
	}

	value, _ := io.ReadAll(jsonFile)
	json.Unmarshal(value, &testSetts)
	return testSetts, nil
}

func getSettings() (DirSettings, error) {
	var Setts DirSettings
	jsonFile, err := os.Open("./settings.json")
	if err != nil {
		return Setts, err
	}

	value, _ := io.ReadAll(jsonFile)
	json.Unmarshal(value, &Setts)
	return Setts, nil
}

func finLocationByName(name string) (location, error) {

	for i := 0; i < len(testLocation); i++ {
		if name == testLocation[i].name {
			return testLocation[i], nil
		}
	}

	var loc location
	return loc, &customError{Message: "Invalid name"}
}

func SettingsAndPrepareTestFiles(nameOfLogFile string) DirSettings {
	settings, err := getSettings()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(2)
	}

	DocFilesFormat := strings.Fields(settings.DocumentFiles)
	PicFilesFormat := strings.Fields(settings.PictureFiles)

	var files []*os.File
	err = filepath.Walk("../TestFilesAll", func(path string, info fs.FileInfo, err error) error {
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

	/*for i := 0; i < len(files); i++ {
		fmt.Println(files[i].Name())
	}*/

	for i := 0; i < len(locationForTestFolder); i++ {
		helpers.CreateMultipleTestFiles(locationForTestFolder[i], nameOfLogFile, files, settings.DirNumber, settings.DirNumberFiles)
	}

	return settings
}

func getWhichTestToExecute(nameOfLogFile string) []location {
	tests, err := getTestSettings()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(1)
	}

	var whichTests []location
	for i := 1; i < len(tests); i++ {
		if tests[i].IsEnabled {
			if tests[i].Name[:3] == "Win" || tests[i].Name[:3] == "Lin" {
				if strings.ToLower(tests[i].Name[:3]) == runtime.GOOS[:3] {
					loc, err := finLocationByName(tests[i].Name[3:])
					if err != nil {
						fmt.Println(err.Error())
						os.Exit(1)
					}
					whichTests = append(whichTests, loc)
				}

			} else {
				loc, err := finLocationByName(tests[i].Name[3:])
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
				whichTests = append(whichTests, loc)
			}

		}
	}

	return whichTests
}

func main() {
	nameOfLogFile := helpers.CreateLogFileIfItDoesNotExist("./", "main")

	whichTests := getWhichTestToExecute(nameOfLogFile)

	settings := SettingsAndPrepareTestFiles(nameOfLogFile)

	obfuscate.PrepareEveryTestObfuscated(nameOfLogFile)

	var correctTests, incorrectTests []string
	for i := 0; i < len(whichTests); i++ {
		fmt.Println("Starting test: " + whichTests[i].name)
		if whichTests[i].name == "EncryptDecryptDirRecursive" || whichTests[i].name == "EncryptDecryptDirRecursivePartially" {
			cmd := exec.Command(whichTests[i].path+whichTests[i].nameOfFile, settings.ToEncrypt, settings.ToDecrypt)
			cmd.Dir = whichTests[i].path
			err := cmd.Run()

			if err == nil {
				correctTests = append(correctTests, whichTests[i].nameOfFile)
			} else {
				helpers.WriteLog(nameOfLogFile, err.Error()+" from "+whichTests[i].name, 1)
				incorrectTests = append(incorrectTests, whichTests[i].nameOfFile)
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
		} else if whichTests[i].name == "PrivilegeEscalation" {
			var pid string
			fmt.Println("Choose which pid to duplicate:")
			fmt.Scanln(&pid)
			cmd := exec.Command(whichTests[i].path+whichTests[i].nameOfFile, pid)
			cmd.Dir = whichTests[i].path
			err := cmd.Run()

			if err.Error() == whichTests[i].expectedResult {
				correctTests = append(correctTests, whichTests[i].name)
			} else {
				helpers.WriteLog(nameOfLogFile, err.Error()+" from "+whichTests[i].name, 1)
				incorrectTests = append(incorrectTests, whichTests[i].name)
			}
		} else if whichTests[i].name == "ProcessHollowing" {
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
	fmt.Printf("Out of %d tests, %d are/is correct, %d are/is incorrect", len(whichTests), len(correctTests), len(incorrectTests))
	for j := 0; j < len(correctTests); j++ {
		fmt.Printf("\nTest %s is correct!", correctTests[j])
	}
	for j := 0; j < len(incorrectTests); j++ {
		fmt.Printf("\nTest %s is incorrect!", incorrectTests[j])
	}

	for i := 0; i < len(locationForTestFolder); i++ {
		err := os.RemoveAll(locationForTestFolder[i] + "/testFilesParent")
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	}

}
