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
	"time"
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

type Result struct {
	Name         string
	Timestamp    string
	ExitStatus   string
	ErrorMessage string
	IsTrue       bool
}

var testLocationWin = []location{
	location{"ArchiveFiles", "../ArchiveFiles/", "ArchiveFiles.exe", "nil"},
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
	location{"OSCredentialDump", "../OSCredDump/", "CredDump.exe", "exit status 1"},
	location{"DefenseEnvasion", "../DefEnv/", "DefEnv.exe", "exit status 1"}}

var testLocationLin = []location{
	location{"ArchiveFiles", "../ArchiveFiles/", "ArchiveFiles", "nil"},
	location{"Eicar", "../Eicar/", "Eic", "exit status 1"},
	location{"EncryptDecryptDirRecursive", "../EncryptDecryptDirRecursive/", "EncryptDecryptDirRecursive", "nil"},
	location{"EncryptDecryptDirRecursivePartially", "../EncryptDecryptDirRecursivePartially/", "EncryptDecryptDirRecursivePartially", "nil"},
	location{"GetSystemInformation", "../GetSysInfo/", "GetSysInfo", "nil"},
	location{"SecureDeleteFiles", "../SecureDeleteFiles/", "SecureDeleteFiles", "nil"},
	location{"CronJob", "../CronJob/", "Cron", "nil"},
	location{"RansomwareNoteDeploy", "../RanNote/", "RanNote", "nil"},
	location{"ExfiltrationOfFiles", "../UploadFiles/", "Exfilt", "nil"}}

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

	if runtime.GOOS == "windows" {
		for i := 0; i < len(testLocationWin); i++ {
			if name == testLocationWin[i].name {
				return testLocationWin[i], nil
			}
		}
	} else {
		for i := 0; i < len(testLocationLin); i++ {
			if name == testLocationLin[i].name {
				return testLocationLin[i], nil
			}
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
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "main")
		os.Exit(1)
	}

	for i := 0; i < len(locationForTestFolder); i++ {
		helpers.CreateMultipleTestFiles(locationForTestFolder[i], nameOfLogFile, files, settings.DirNumber, settings.DirNumberFiles, "main")
	}

}

func getWhichTestToExecuteAndSettings(nameOfLogFile string) ([]location, Setts) {
	config, err := getTestConfig()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "main")
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

func getResults(nameOfLogFile string, testResults []Result) {
	os.Remove("./results.html")
	f, err := os.Create("./results.html")
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "main")
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
		"<tr>" + "<th>Name Of Test</th>" + "<th>Did it recieve expected exit status</th>" + "<th>Time</th>" + "<th>Recieved exit status</th>" + "<th>Recieved Error</th>" + "</tr>"
	for i := 0; i < len(testResults); i++ {
		fmt.Println("Test:")
		fmt.Println(testResults[i].Name)
		fmt.Println("Did it recieve expected exit status:")
		isTrueString := "no"
		if testResults[i].IsTrue {
			isTrueString = "yes"
		}
		fmt.Println(isTrueString)
		fmt.Println("Timestamp:")
		fmt.Println(testResults[i].Timestamp)
		fmt.Println("Exit status:")
		fmt.Println(testResults[i].ExitStatus)
		fmt.Println("Error message:")
		fmt.Println(testResults[i].ErrorMessage)
		fmt.Println("")
		if testResults[i].IsTrue {
			content += "<tr>" + "<td style='background-color: lightgreen'>" + testResults[i].Name + "</td>" + "<td style='background-color: lightgreen' >Yes</td>" + "<td style='background-color: lightgreen' >" + testResults[i].Timestamp + "</td>" + "<td style='background-color: lightgreen' >" + testResults[i].ExitStatus + "</td>" + "<td style='background-color: lightgreen' >" + testResults[i].ErrorMessage + "</td>" + "<tr>"
		} else {
			content += "<tr>" + "<td style='background-color: LightCoral'>" + testResults[i].Name + "</td>" + "<td style='background-color: LightCoral' >No</td>" + "<td style='background-color: LightCoral' >" + testResults[i].Timestamp + "</td>" + "<td style='background-color: LightCoral' >" + testResults[i].ExitStatus + "</td>" + "<td style='background-color: LightCoral' >" + testResults[i].ErrorMessage + "</td>" + "<tr>"
		}
	}

	content += "</table>" + "</body>" + "</html>"
	f.WriteString(content)
}

func getLastError(name string, logfile string) string {
	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("file openning")
		os.Exit(100)
	}

	data, err := os.ReadFile(f.Name())
	if err != nil {
		fmt.Println("file reading")
		os.Exit(101)
	}

	var jsonValues []helpers.Log

	err = json.Unmarshal(data, &jsonValues)
	if err != nil {
		fmt.Println("Unmarshall")
		os.Exit(102)
	}

	for i := len(jsonValues) - 1; i >= 0; i-- {
		if jsonValues[i].TypeOfLog == "ERROR" && jsonValues[i].Test == name {
			return jsonValues[i].Line
		}
	}

	return ""
}

func main() {
	nameOfLogFile := helpers.CreateLogFileIfItDoesNotExist("./", "main", "main")

	whichTests, settings := getWhichTestToExecuteAndSettings("tests.json")
	prepareTestFiles(nameOfLogFile, settings)
	//obfuscate.PrepareEveryTestObfuscated(nameOfLogFile)

	var testResults []Result
	for i := 0; i < len(whichTests); i++ {
		fmt.Println("Starting test: " + whichTests[i].name)
		if whichTests[i].name == "EncryptDecryptDirRecursive" || whichTests[i].name == "EncryptDecryptDirRecursivePartially" {
			cmd := exec.Command(whichTests[i].path+whichTests[i].nameOfFile, settings.ToEncrypt, settings.ToDecrypt, nameOfLogFile)
			cmd.Dir = whichTests[i].path
			err := cmd.Run()

			result := Result{whichTests[i].name, time.Now().Format(time.RFC822), "", getLastError(whichTests[i].name, nameOfLogFile), false}
			if err == nil {
				result.IsTrue = true
				result.ExitStatus = "exit status 0"
			} else {
				result.IsTrue = false
				result.ExitStatus = err.Error()
			}

			testResults = append(testResults, result)
		} else if whichTests[i].name == "ServiceCreation" {
			cmd := exec.Command(whichTests[i].path+whichTests[i].nameOfFile, "install", nameOfLogFile)
			cmd.Dir = whichTests[i].path
			err := cmd.Run()

			result := Result{whichTests[i].name, time.Now().Format(time.RFC822), "", getLastError(whichTests[i].name, nameOfLogFile), false}

			if err == nil {
				if whichTests[i].expectedResult == "nil" {
					result.IsTrue = true
				} else {
					result.IsTrue = false
				}

				result.ExitStatus = "exit status 0"

				cmd := exec.Command(whichTests[i].path+whichTests[i].nameOfFile, "uninstall", nameOfLogFile)
				cmd.Dir = whichTests[i].path
				cmd.Run()
			} else {
				if err.Error() == whichTests[i].expectedResult {
					result.IsTrue = true
				} else {
					result.IsTrue = false
				}

				result.ExitStatus = err.Error()
			}

			testResults = append(testResults, result)
		} else if whichTests[i].name == "SecureDeleteFiles" || whichTests[i].name == "ArchiveFiles" {
			cmd := exec.Command(whichTests[i].path+whichTests[i].nameOfFile, strconv.Itoa(settings.TimeToDelayOnSecureDelete), nameOfLogFile)
			cmd.Dir = whichTests[i].path
			err := cmd.Run()

			result := Result{whichTests[i].name, time.Now().Format(time.RFC822), "", getLastError(whichTests[i].name, nameOfLogFile), false}

			if err == nil {
				if whichTests[i].expectedResult == "nil" {
					result.IsTrue = true
				} else {
					result.IsTrue = false
				}

				result.ExitStatus = "exit status 0"
			} else {
				if err.Error() == whichTests[i].expectedResult {
					result.IsTrue = true
				} else {
					result.IsTrue = false
				}
				result.ExitStatus = err.Error()
			}

			testResults = append(testResults, result)
		} else if whichTests[i].name == "ProcessHollowing32" {
			cmd := exec.Command(whichTests[i].path+whichTests[i].nameOfFile, "C:\\windows\\syswow64\\notepad.exe")
			cmd.Dir = whichTests[i].path
			err := cmd.Run()

			result := Result{whichTests[i].name, time.Now().Format(time.RFC822), "", getLastError(whichTests[i].name, nameOfLogFile), false}

			if err == nil {
				result.IsTrue = false
				result.ExitStatus = "exit status 0"
			} else {
				if err.Error() == whichTests[i].expectedResult {
					result.IsTrue = true
				} else {
					result.IsTrue = false
				}
				result.ExitStatus = err.Error()
			}

			testResults = append(testResults, result)
		} else if whichTests[i].name == "ProcessHollowing64" {
			cmd := exec.Command(whichTests[i].path+whichTests[i].nameOfFile, "C:\\windows\\System32\\notepad.exe")
			cmd.Dir = whichTests[i].path
			err := cmd.Run()

			result := Result{whichTests[i].name, time.Now().Format(time.RFC822), "", getLastError(whichTests[i].name, nameOfLogFile), false}

			if err == nil {
				result.IsTrue = false
				result.ExitStatus = "exit status 0"
			} else {
				if err.Error() == whichTests[i].expectedResult {
					result.IsTrue = true
				} else {
					result.IsTrue = false
				}
				result.ExitStatus = err.Error()
			}

			testResults = append(testResults, result)
		} else if whichTests[i].name == "ReverseShell" {
			cmd := exec.Command(whichTests[i].path+whichTests[i].nameOfFile, settings.IPAddrForReverseShell, nameOfLogFile)
			cmd.Dir = whichTests[i].path
			err := cmd.Run()

			result := Result{whichTests[i].name, time.Now().Format(time.RFC822), "", getLastError(whichTests[i].name, nameOfLogFile), false}

			if err == nil {
				result.IsTrue = false
				result.ExitStatus = "exit status 0"
			} else {
				if err.Error() == whichTests[i].expectedResult {
					result.IsTrue = true
				} else {
					result.IsTrue = false
				}
				result.ExitStatus = err.Error()
			}

			testResults = append(testResults, result)
		} else if whichTests[i].name == "ExfiltrationOfFiles" {
			cmd := exec.Command(whichTests[i].path+whichTests[i].nameOfFile, settings.MegaEmail, settings.MegaPassword, "UploadFiles")
			cmd.Dir = whichTests[i].path
			err := cmd.Run()

			result := Result{whichTests[i].name, time.Now().Format(time.RFC822), "", getLastError(whichTests[i].name, nameOfLogFile), false}

			if err == nil {
				result.IsTrue = true
				result.ExitStatus = "exit status 0"
			} else {
				result.IsTrue = false
				result.ExitStatus = err.Error()
			}

			testResults = append(testResults, result)
		} else {
			cmd := exec.Command(whichTests[i].path+whichTests[i].nameOfFile, nameOfLogFile)
			cmd.Dir = whichTests[i].path
			err := cmd.Run()

			result := Result{whichTests[i].name, time.Now().Format(time.RFC822), "", getLastError(whichTests[i].name, nameOfLogFile), false}

			if err == nil {
				if whichTests[i].expectedResult == "nil" {
					result.IsTrue = true
				} else {
					result.IsTrue = false
				}
				result.ExitStatus = "exit status 0"
			} else {
				if err.Error() == whichTests[i].expectedResult {
					result.IsTrue = true
				} else {
					result.IsTrue = false
				}
				result.ExitStatus = err.Error()
			}

			testResults = append(testResults, result)
		}
		fmt.Println("Endinging test: " + whichTests[i].name)
	}

	for i := 0; i < len(whichTests); i++ {
		_, err := os.Open(whichTests[i].path + whichTests[i].nameOfFile)
		if err != nil {
			for j := 0; j < len(testResults); j++ {
				if testResults[j].Name == whichTests[i].name {
					if whichTests[i].expectedResult == "exit status 1" {
						testResults[j].IsTrue = true
					} else {
						testResults[j].IsTrue = false
					}
				}
			}
		}
	}

	getResults(nameOfLogFile, testResults)
	fmt.Println("You can check the results!")

	for i := 0; i < len(locationForTestFolder); i++ {
		err := helpers.RemoveTestFilesIfExists(locationForTestFolder[i])
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1, "main")
		}
	}

}
