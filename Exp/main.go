package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Please input payload and target path!")
		os.Exit(2)
	}
	payloadPath := os.Args[1]
	targetPath := os.Args[2]
	if HollowProcess(payloadPath, targetPath) {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

//a
