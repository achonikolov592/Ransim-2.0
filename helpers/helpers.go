package helpers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Log struct {
	TypeOfLog string `json:"type"`
	Line      string `json:"line"`
	Test      string `json:"test"`
	Time      string `json:"time"`
}

func AppendAtTheEnd(file string, logfile string, nameOfTest string) {
	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		WriteLog(logfile, err.Error(), 1, nameOfTest)
	}

	_, err = f.Write([]byte("//a"))
	if err != nil {
		WriteLog(logfile, err.Error(), 1, nameOfTest)
	}
}

func CreateTestFiles(dir string, logfile string, nameOfTest string) {
	WriteLog(logfile, "Start creating test files", 2, nameOfTest)

	if _, err := os.Stat(dir + "testfiles"); err != nil {
		err := os.Mkdir(dir+"testfiles", 0777)
		if err != nil {
			WriteLog(logfile, err.Error(), 1, nameOfTest)
			os.Exit(1)
		}
	}

	if _, err := os.Stat(dir + "testfiles/sub"); err != nil {
		err := os.Mkdir(dir+"testfiles/sub", 0777)
		if err != nil {
			WriteLog(logfile, err.Error(), 1, nameOfTest)
			os.Exit(1)
		}
	}

	f, err := os.Create(dir + "testfiles/c.txt")
	if err != nil {
		WriteLog(logfile, err.Error(), 1, nameOfTest)
		os.Exit(2)
	}

	_, err = f.WriteString("asdfghjk")
	if err != nil {
		WriteLog(logfile, err.Error(), 1, nameOfTest)
		os.Exit(2)
	}

	f, err = os.Create(dir + "testfiles/sub/b.txt")
	if err != nil {
		WriteLog(logfile, err.Error(), 1, nameOfTest)
		os.Exit(2)
	}

	_, err = f.WriteString("asdfghjk")
	if err != nil {
		WriteLog(logfile, err.Error(), 1, nameOfTest)
		os.Exit(2)
	}
}

func CreateMultipleTestFiles(dir string, logfile string, files []*os.File, numberOfDirs, numberOfFIles int, nameOfTest string) {
	if _, err := os.Stat(dir + "/testFilesParent"); err != nil {
		err := os.Mkdir(dir+"/testFilesParent", 0777)
		if err != nil {
			WriteLog(logfile, err.Error(), 1, nameOfTest)
		}
	}

	for i := 0; i < numberOfDirs; i++ {
		testFileDir := dir + "/testFilesParent/testfiles" + strconv.Itoa(i) + "/"
		if _, err := os.Stat(testFileDir); err != nil {
			err := os.Mkdir(testFileDir, 0777)
			if err != nil {
				WriteLog(logfile, err.Error(), 1, nameOfTest)
			}
		}

		for j := 0; j < len(files); j++ {
			fileContent, err := io.ReadAll(files[j])
			if err != nil {
				WriteLog(logfile, err.Error(), 1, nameOfTest)
			}
			files[j].Seek(0, io.SeekStart)
			for iof := 0; iof < numberOfFIles; iof++ {
				f, err := os.Create(testFileDir + strconv.Itoa(iof) + files[j].Name()[16:])
				if err != nil {
					WriteLog(logfile, err.Error(), 1, nameOfTest)
				}
				_, err = f.WriteString(string(fileContent))
				if err != nil {
					WriteLog(logfile, err.Error(), 1, nameOfTest)
				}
			}
		}

		testFileSubDir := testFileDir + "sub/"
		if _, err := os.Stat(testFileSubDir); err != nil {
			err := os.Mkdir(testFileSubDir, 0777)
			if err != nil {
				WriteLog(logfile, err.Error(), 1, nameOfTest)
			}
		}

		for j := 0; j < len(files); j++ {
			fileContent, err := io.ReadAll(files[j])
			if err != nil {
				WriteLog(logfile, err.Error(), 1, nameOfTest)
			}
			files[j].Seek(0, io.SeekStart)
			for iof := 0; iof < numberOfFIles; iof++ {
				f, err := os.Create(testFileSubDir + strconv.Itoa(iof) + files[j].Name()[16:])
				if err != nil {
					WriteLog(logfile, err.Error(), 1, nameOfTest)
				}
				_, err = f.WriteString(string(fileContent))
				if err != nil {
					WriteLog(logfile, err.Error(), 1, nameOfTest)
				}
			}
		}
	}

}

func CreateLogFileIfItDoesNotExist(dir string, name string, nameOfTest string) string {
	file, err := os.Create(dir + name + ".json")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	file.WriteString("[]")

	filepath, err := filepath.Abs(file.Name())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return filepath

}

func RemoveTestFilesIfExists(dir string) error {
	err := os.RemoveAll(dir + "testfiles")
	if err != nil {
		return err
	}
	err = os.RemoveAll(dir + "testFilesParent")

	return err
}

func WriteLog(logfile string, line string, opt int, nameOfTest string) {
	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("file")
		os.Exit(100)
	}

	data, err := os.ReadFile(f.Name())
	if err != nil {
		fmt.Println("oops")
		os.Exit(101)
	}

	var jsonValues []Log

	err = json.Unmarshal(data, &jsonValues)
	if err != nil {
		fmt.Println("Unmarshall")
		os.Exit(102)
	}

	var valueToWrite Log

	valueToWrite.Line = line

	if opt == 1 { //err
		valueToWrite.TypeOfLog = "ERROR"
	} else { //info
		valueToWrite.TypeOfLog = "INFO"
	}

	valueToWrite.Time = time.Now().Format(time.RFC822)

	valueToWrite.Test = nameOfTest

	jsonValues = append(jsonValues, valueToWrite)

	newData, err := json.MarshalIndent(jsonValues, "", "	")
	if err != nil {
		fmt.Println("Marshall")
		os.Exit(103)
	}

	err = os.WriteFile(f.Name(), newData, 0666)
	if err != nil {
		fmt.Println("Write")
		os.Exit(104)
	}
	//fmt.Println(line)

}

func CreateNormalLogFileIfItDoesNotExist(dir string, name string) string {

	i, err := os.Stat(dir + name + ".log")

	if err != nil {
		_, err := os.Create(dir + name + ".log")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		return dir + name + ".log"
	}

	return i.Name()

}

func WriteNormalLog(logfile string, line string) {
	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println(f)
		os.Exit(100)
	}

	var stringToWrite = line + "\n"
	_, err = f.WriteString(stringToWrite)
	if err != nil {
		fmt.Println(err)
		os.Exit(101)
	}
}
