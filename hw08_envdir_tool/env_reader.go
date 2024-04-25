package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envs := make(Environment)

	for _, file := range files {
		filName := file.Name()

		if strings.Contains(filName, "=") {
			continue
		}

		filePath := filepath.Join(dir, filName)

		f, err := os.Open(filePath)
		if err != nil {
			log.Printf("Failed to open file %s\n", filePath)
			continue
		}

		info, err := f.Stat()
		if err != nil {
			if err := f.Close(); err != nil {
				log.Printf("Failed to close file %s\n", f.Name())
				return nil, err
			}
			fmt.Println(err)
			return nil, err
		}

		if info.Size() == 0 {
			envs[filName] = EnvValue{"", true}
		} else {
			scanner := bufio.NewScanner(f)
			scanner.Split(bufio.ScanLines)

			if scanner.Scan() {
				str := scanner.Text()
				str = strings.TrimRight(str, " \t")
				str = string(bytes.ReplaceAll([]byte(str), []byte("\x00"), []byte("\n")))

				envs[filName] = EnvValue{str, false}
			}
		}

		if err = f.Close(); err != nil {
			log.Printf("Failed to close file %s\n", f.Name())
			return nil, err
		}
	}

	return envs, nil
}
