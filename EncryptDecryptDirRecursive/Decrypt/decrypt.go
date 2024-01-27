package decrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"helpers"
	"os"
	"path/filepath"
)

func decrypt(dirToDecrypt string, c cipher.AEAD, nameOfLogFile string) {
	var filesInDir []string
	err := filepath.Walk(dirToDecrypt, func(path string, info os.FileInfo, err error) error {
		if path != dirToDecrypt {
			filesInDir = append(filesInDir, path)
		}
		return nil
	})
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(4)
	}

	for _, file := range filesInDir {
		info, err := os.Stat(file)
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1)
			os.Exit(5)
		}

		if !(info.IsDir()) {
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
				os.Exit(6)
			}
			encryptedText, err := os.ReadFile(file)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
				os.Exit(7)
			}

			non, text := encryptedText[:c.NonceSize()], encryptedText[c.NonceSize():]

			decb, err := c.Open(nil, non, text, nil)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
				os.Exit(8)
			}
			err = os.WriteFile(file, decb, 0666)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
				os.Exit(9)
			}

		} else {
			go decrypt(info.Name(), c, nameOfLogFile)
		}
	}

}

func DecryptDir(dirToDecrypt string, nameOfEncryptionInfoFile string, nameOfLogFile string) {
	file, _ := os.ReadFile(nameOfEncryptionInfoFile)
	key, _ := hex.DecodeString(string(file[0:64]))

	block, err := aes.NewCipher(key)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(2)
	}
	c, err := cipher.NewGCM(block)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(3)
	}
	decrypt(dirToDecrypt, c, nameOfLogFile)
}
