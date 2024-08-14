package zipdir

import (
	"RRA/SecureDeleteFiles/SecureDeleteFile"
	"archive/zip"
	"helpers"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func ZipDir(dir string, nameOfLogFile string, timeToDelay int, nameOfTest string) {
	var files []string

	errFileWalk := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	if errFileWalk != nil {
		helpers.WriteLog(nameOfLogFile, errFileWalk.Error(), 1, "Zip Dir from"+nameOfTest)
		os.Exit(2)
	}

	archive, err := os.Create("./archive.zip")
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, "Zip Dir from"+nameOfTest)
		os.Exit(3)
	}
	defer archive.Close()
	zipWriting := zip.NewWriter(archive)

	for _, file := range files {
		archivedFile, errOfOpeningFile := os.Open(file)
		if errOfOpeningFile != nil {
			helpers.WriteLog(nameOfLogFile, errOfOpeningFile.Error(), 1, "Zip Dir from"+nameOfTest)
			os.Exit(4)
		}

		fileInfo, _ := os.Stat(file)

		if !fileInfo.IsDir() {
			w1, err := zipWriting.Create(file)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1, "Zip Dir from"+nameOfTest)
				os.Exit(5)
			}
			if _, err := io.Copy(w1, archivedFile); err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1, "Zip Dir from"+nameOfTest)
				os.Exit(6)
			}
			SecureDeleteFile.SecureDelete(file, nameOfLogFile, timeToDelay, nameOfTest)
		}

		archivedFile.Close()
	}

	zipWriting.Close()
	os.RemoveAll(dir)
}
