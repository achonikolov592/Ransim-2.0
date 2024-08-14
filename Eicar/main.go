package main

import (
	"fmt"
	"helpers"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	var nameOfLogFile string
	if len(os.Args) == 2 {
		nameOfLogFile = os.Args[1]
	} else {
		nameOfLogFile = helpers.CreateLogFileIfItDoesNotExist("./", "Eicar", "Eicar")
	}

	helpers.WriteLog(nameOfLogFile, "Starting test: EicarTest", 2, "Eicar")

	out, err := os.Create("./out.txt")
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "Eicar")
		os.Exit(2)
	}
	defer out.Close()

	resp, err := http.Get("https://www.eicar.org/download/eicar-com-2/?wpdmdl=8842&refresh=656a05f4d5e5d1701447156")
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "Eicar")
		os.Exit(3)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "Eicar")
		os.Exit(4)
	}

	time.Sleep(10 * time.Second)

	_, err = os.Open("./out.txt")
	fmt.Println(err)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "Eicar")
		os.Exit(1)
	}

	helpers.WriteLog(nameOfLogFile, "Ending test: EicarTest", 2, "Eicar")
	os.Exit(0)
}

//as
