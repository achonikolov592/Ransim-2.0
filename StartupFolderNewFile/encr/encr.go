package main

import (
	"RRA/EncryptDecryptDirRecursive/encrypt"
	"os"
	"path/filepath"
)

func main() {
	ex, _ := os.Executable()
	exPath, _ := filepath.Abs(ex)
	encrypt.EncryptDir(exPath+"/../testFilesParent", exPath+"/../startup.log", exPath+"/../StartupFolderNewFile/EncryptionInfo.log")
}
