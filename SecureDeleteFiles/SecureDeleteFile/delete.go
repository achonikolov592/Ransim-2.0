package SecureDeleteFile

import (
	"helpers"
	"os"
	"time"
)

func SecureDelete(filename string, nameOfLogFile string, timeToDelay int) {
	file, err := os.OpenFile(filename, os.O_WRONLY, 0666)

	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(5)
	}

	stats, err := file.Stat()

	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(6)
	}

	sizeOfFile := stats.Size()
	chunk := int64(1 * (1 << 18))

	numberOfParts := (sizeOfFile / chunk) + 1
	position := int64(0)

	for i := int64(0); i < numberOfParts; i++ {
		var howMuchToChange int64
		if sizeOfFile-(i+1)*chunk < 0 {
			howMuchToChange = sizeOfFile - i*chunk
		} else {
			howMuchToChange = chunk
		}

		zeros := make([]byte, howMuchToChange)

		_, err = file.WriteAt([]byte(zeros), position)
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1)
			os.Exit(7)
		}

		position += howMuchToChange
	}

	file.Close()

	time.Sleep(time.Second * time.Duration(timeToDelay))
	if err = os.Remove(filename); err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(8)
	}

}
