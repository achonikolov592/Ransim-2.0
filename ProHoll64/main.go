package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please input target path!")
		os.Exit(1)
	}
	targetPath := os.Args[1]
	HollowProcess("./processToHollow/Process64.exe", targetPath)
}

//a
