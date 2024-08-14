package main

import (
	"helpers"
	"io"
	"net/http"
	"os"
)

var nameOfLogFile string

func main() {
	if len(os.Args) == 2 {
		nameOfLogFile = os.Args[1]
	} else {
		nameOfLogFile = helpers.CreateLogFileIfItDoesNotExist("./", "DownloadMalicious", "DownloadMalicious")
	}

	helpers.WriteLog(nameOfLogFile, "Starting test: DownloadingMaliciousPayload", 1, "DownloadMalicious")

	out, err := os.Create("./out1.zip")
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "DownloadMalicious")
		os.Exit(2)
	}
	defer os.Remove("./out1.zip")

	resp, err := http.Get("https://github.com/kh4sh3i/Ransomware-Samples/raw/main/Petya/Ransomware.Petya.zip")
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "DownloadMalicious")
		os.Exit(3)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "DownloadMalicious")
		os.Exit(4)
	}

	_, err = os.Open("./out1.zip")
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "DownloadMalicious")
		os.Exit(1)
	}

	helpers.WriteLog(nameOfLogFile, "Ending test: DownloadingMaliciousPayload", 1, "DownloadMalicious")

	os.Exit(0)
}

//a
