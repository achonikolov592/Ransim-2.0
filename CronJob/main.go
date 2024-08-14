package main

import (
	"RRA/EncryptDecryptDirRecursive/encrypt"
	"helpers"
	"os"

	"github.com/robfig/cron"
)

func main() {
	var nameOfLogFile string
	if len(os.Args) == 2 {
		nameOfLogFile = os.Args[1]
	} else {
		nameOfLogFile = helpers.CreateLogFileIfItDoesNotExist("./", "CronJob", "CronJob")
	}

	nameOfEncryptionInfoFile := helpers.CreateNormalLogFileIfItDoesNotExist("./", "EncryptionInfo")
	helpers.WriteLog(nameOfLogFile, "Starting test: CronJob", 2, "CronJob")
	c := cron.New()
	i := 0
	sig := make(chan int)
	c.AddFunc("@every 10s", func() {
		if i < 6 {
			helpers.RemoveTestFilesIfExists("./")
			helpers.CreateTestFiles("./", nameOfLogFile, "CronJob")
			encrypt.EncryptDir("./testfiles", nameOfLogFile, nameOfEncryptionInfoFile, "CronJob")
			helpers.WriteLog(nameOfLogFile, "Encrypted Files", 2, "CronJob")
			i++
		} else {
			c.Stop()
			helpers.WriteLog(nameOfLogFile, "Ending test: CronJob", 2, "CronJob")
			os.Exit(0)
		}

	})

	c.Start()

	<-sig

}
