package main

import (
	"io"
	"log"
	"os"
	"os/exec"
)

const errorcode = 1

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return errorcode
	}

	name := cmd[0]
	command := exec.Command(name, cmd[1:]...)

	for s, t := range env {
		os.Unsetenv(s)
		if t.Value != "" {
			os.Setenv(s, t.Value)
		}
	}

	stdin, err := command.StdinPipe()
	if err != nil {
		log.Println(err)
	}
	go func() {
		defer stdin.Close()
		io.Copy(stdin, os.Stdin)
	}()
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err = command.Start(); err != nil {
		log.Println(err)
	}

	return returnCode
}
