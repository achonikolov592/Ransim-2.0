package main

import (
	"helpers"
	"io"
	"net/http"
	"os"
)

var nameOfLogFile string

func main() {
	nameOfLogFile = helpers.CreateLogFileIfItDoesNotExist("./", "eicar")

	helpers.WriteLog(nameOfLogFile, "Starting test: EicarTest", 2)

	out, err := os.Create("./out.txt")
	if err != nil {
		os.Exit(2)
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	}
	defer out.Close()

	resp, err := http.Get("https://www.eicar.org/download/eicar-com-2/?wpdmdl=8842&refresh=656a05f4d5e5d1701447156")
	if err != nil {
		os.Exit(3)
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(4)
	}

	_, err = os.Open("./out1.txt")
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(1)
	}

	helpers.WriteLog(nameOfLogFile, "Ending test: EicarTest", 2)
	os.Exit(0)
}

//as
