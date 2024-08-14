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

func EncryptFilesInDir(dirToEncrypt string, nameOfLogFile string, nameOfEncryptionInfoFile string, nameOfTest string) {
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

		f, err := os.Open(file)
		if err != nil {
			helpers.WriteLog(nameOfLogFile, err.Error(), 1, nameOfTest)
			os.Exit(7)
		}

		if !(info.IsDir()) {
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1, nameOfTest)
				os.Exit(8)
			}

			howMuchToEncrypt := info.Size() / int64(8)
			whereToStart := int64(howMuchToEncrypt * mathrand.Int63n(8))
			filetext, err := os.ReadFile(file)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1, nameOfTest)
				os.Exit(9)
			}
			textToEncrypt := make([]byte, howMuchToEncrypt)
			_, err = f.ReadAt(textToEncrypt, whereToStart)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1, nameOfTest)
				os.Exit(10)
			}

			nonce := make([]byte, c.NonceSize())
			if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1, nameOfTest)
				os.Exit(11)
			}
			encryptedText := c.Seal(nonce, nonce, textToEncrypt, nil)

			helpers.WriteNormalLog(nameOfEncryptionInfoFile, "Where it started to encrypt: "+strconv.Itoa(int(whereToStart)))
			helpers.WriteNormalLog(nameOfEncryptionInfoFile, "How much bytes is product: "+strconv.Itoa(len(encryptedText)))

			var finaltext []byte
			finaltext = append(finaltext, filetext[:whereToStart]...)
			finaltext = append(finaltext, encryptedText...)
			finaltext = append(finaltext, filetext[whereToStart+howMuchToEncrypt:]...)
			err = os.WriteFile(file, finaltext, 0666)
			if err != nil {
				helpers.WriteLog(nameOfLogFile, err.Error(), 1, nameOfTest)
				os.Exit(12)
			}

			helpers.WriteLog(nameOfLogFile, "Encrypted: "+file, 2, nameOfTest)

		}
	}

} //a
