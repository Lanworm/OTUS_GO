package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const testDir = "./testdata/env"

func TestReadDir(t *testing.T) {
	expected := Environment{
		"BAR":   EnvValue{"bar", false},
		"EMPTY": EnvValue{"", false},
		"FOO":   EnvValue{"   foo\nwith new line", false},
		"HELLO": EnvValue{"\"hello\"", false},
		"UNSET": EnvValue{"", true},
	}

	t.Run("default test", func(t *testing.T) {
		envs, err := ReadDir(testDir)

		require.Nil(t, err)
		require.Equal(t, expected, envs)
	})

	t.Run("dir not exists", func(t *testing.T) {
		_, err := ReadDir("test_test")
		require.NotEqual(t, nil, err)
	})

	t.Run("filename contains character '='", func(t *testing.T) {
		envTest, err := os.CreateTemp(testDir, "out=*")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer os.Remove(envTest.Name())

		envs, err := ReadDir(testDir)

		require.Nil(t, err)
		require.Equal(t, expected, envs)
	})
}
