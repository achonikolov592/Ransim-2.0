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

func encrypt(dirToEncrypt string, c cipher.AEAD, nameOfLogFile string) {
	var filesInDir []string
	err := filepath.Walk(dirToEncrypt, func(path string, info os.FileInfo, err error) error {
		if path != dirToEncrypt {
			filesInDir = append(filesInDir, path)
		}
		return nil
	})
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(5)
	}

	for _, file := range filesInDir {
		info, err := os.Stat(file)
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1)
			os.Exit(6)
		}

		if !(info.IsDir()) {
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
				os.Exit(7)
			}
			encryptedText, err := os.ReadFile(file)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
				os.Exit(8)
			}

			nonce := make([]byte, c.NonceSize())
			if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
				os.Exit(9)
			}
			encryptedText = c.Seal(nonce, nonce, encryptedText, nil)

			err = os.WriteFile(file, encryptedText, 0666)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
				os.Exit(10)
			}

			helpers.WriteLog(nameOfLogFile, "Encrypted: "+file, 2)

		} else {
			encrypt(info.Name(), c, nameOfLogFile)
		}
	}

}

func EncryptDir(dirToEncrypt string, nameOfLogFile string, nameOfEncryptionInfoFile string) {

	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(2)
	}

	helpers.WriteLog(nameOfEncryptionInfoFile, hex.EncodeToString(key), 0)

	block, err := aes.NewCipher(key)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(3)
	}

	c, err := cipher.NewGCM(block)
	if err != nil {
		helpers.WriteLog(nameOfLogFile, err.Error(), 1)
		os.Exit(4)
	}

	encrypt(dirToEncrypt, c, nameOfLogFile)
}
