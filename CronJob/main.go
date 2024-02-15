package main

import (
	"RRA/EncryptDecryptDirRecursive/encrypt"
	"helpers"
	"os"

	"github.com/robfig/cron"
)

func main() {
	nameOfLogFile := helpers.CreateLogFileIfItDoesNotExist("./", "CronJob")
	nameOfEncryptionInfoFile := helpers.CreateLogFileIfItDoesNotExist("./", "EncryptionInfo")
	helpers.WriteLog(nameOfLogFile, "Starting CronJob test.", 2)
	c := cron.New()
	i := 0
	sig := make(chan int)
	c.AddFunc("@every 10s", func() {
		if i < 6 {
			helpers.RemoveTestFilesIfExists("./")
			helpers.CreateTestFiles("./", nameOfLogFile)
			encrypt.EncryptDir("./testfiles", nameOfLogFile, nameOfEncryptionInfoFile)
			helpers.WriteLog(nameOfLogFile, "Encrypted Files", 2)
			i++
		} else {
			c.Stop()
			helpers.WriteLog(nameOfLogFile, "Ending CronJob test.", 2)
			os.Exit(0)
		}

	})

	c.Start()

	<-sig

}
