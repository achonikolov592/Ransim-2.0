package obfuscate

import (
	"fmt"
	"helpers"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func PrepareEveryTestObfuscated(nameOfLogFile string) {
	err := filepath.Walk("../", func(path string, info os.FileInfo, err error) error {
		if path != "../" {
			_, err := os.Stat(path)
			if info.IsDir() && !strings.Contains(path, "processToHollow") && !strings.Contains(path, "TestFiles") && !strings.Contains(path, "DLLs") && !strings.Contains(path, "mainAdmin") && !strings.Contains(path, "main") && !strings.Contains(path, ".git") && !strings.Contains(path, "testfiles") && !strings.Contains(path, "testFilesParent") {
				fmt.Println("Building: " + path)
				//cmd := exec.Command("garble", "-literals", "build", ".")
				cmd := exec.Command("go", "build", ".")
				cmd.Dir = path
				err = cmd.Run()
				if err != nil {
					helpers.WriteLog(nameOfLogFile, err.Error()+" at file "+path, 1)
				}
			}
		}
		return nil
	})

	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
	}
}
