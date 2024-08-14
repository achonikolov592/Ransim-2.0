package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"helpers"
	"io"
	"os"
	"path/filepath"
)

func EncryptDir(dirToEncrypt string, nameOfLogFile string, nameOfEncryptionInfoFile string, nameOfTest string) {

	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, nameOfTest)
		os.Exit(2)
	}

	helpers.WriteNormalLog(nameOfEncryptionInfoFile, hex.EncodeToString(key))

	block, err := aes.NewCipher(key)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, nameOfTest)
		os.Exit(3)
	}

	c, err := cipher.NewGCM(block)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, nameOfTest)
		os.Exit(4)
	}

	var filesInDir []string
	err = filepath.Walk(dirToEncrypt, func(path string, info os.FileInfo, err error) error {
		if path != dirToEncrypt {
			filesInDir = append(filesInDir, path)
		}
		return nil
	})
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1, nameOfTest)
		os.Exit(5)
	}

	for _, file := range filesInDir {
		info, err := os.Stat(file)
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1, nameOfTest)
			os.Exit(6)
		}

		if !(info.IsDir()) {
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1, nameOfTest)
				os.Exit(7)
			}
			encryptedText, err := os.ReadFile(file)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1, nameOfTest)
				os.Exit(8)
			}

			nonce := make([]byte, c.NonceSize())
			if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1, nameOfTest)
				os.Exit(9)
			}
			encryptedText = c.Seal(nonce, nonce, encryptedText, nil)

			err = os.WriteFile(file, encryptedText, 0666)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1, nameOfTest)
				os.Exit(10)
			}

			helpers.WriteLog(nameOfLogFile, "Encrypted: "+file, 2, "Encrypt Dir from "+nameOfTest)

		}
	}
}
