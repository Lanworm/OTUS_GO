package main

import (
	"fmt"
	"os"
)

func main() {
	// Place your code here.
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Expected at least 2 arguments")
		return
	}

	dir := args[1]
	command := args[2:]

	envInDir, err := ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		return
	}

	os.Exit(RunCmd(command, envInDir))
}
