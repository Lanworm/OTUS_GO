package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		srcFileName := "testdata/input.txt"
		dstFileName := "testdata/testoutput.txt"
		tests := []struct {
			offset int64
			limit  int64
		}{
			{0, 0},
			{0, 10},
			{0, 1000},
			{0, 10000},
			{100, 1000},
			{6000, 1000},
		}
		for _, tt := range tests {
			t.Run(fmt.Sprintf("offset=%d_limit=%d", tt.offset, tt.limit), func(t *testing.T) {
				err := Copy(srcFileName, dstFileName, tt.offset, tt.limit)
				require.NoError(t, err)

				os.Remove(dstFileName)
			})
		}
	})

	t.Run("errors: negative offset or limit", func(t *testing.T) {
		tests := []struct {
			name   string
			offset int64
			limit  int64
		}{
			{"negative offset", -1, 0},
			{"negative limit", 0, -1},
			{"offset exceeds file size", 10_000, 0},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := Copy("testdata/input.txt", "/tmp/out.txt", tt.offset, tt.limit)
				require.Error(t, err)
			})
		}
	})

	t.Run("errors: empty source or destination", func(t *testing.T) {
		t.Run("empty source name", func(t *testing.T) {
			err := Copy("", "/tmp/out.txt", 0, 0)
			require.Error(t, err)
		})

		t.Run("empty destination name", func(t *testing.T) {
			err := Copy("testdata/input.txt", "", 0, 0)
			require.Error(t, err)
		})

		t.Run("source file doesn't exist", func(t *testing.T) {
			err := Copy("testdata/input1.txt", "/tmp/out.txt", 0, 0)
			require.Error(t, err)
		})

		t.Run("no permission to read source", func(t *testing.T) {
			os.Chmod("testdata/no_permissios.txt", 0o020)

			err := Copy("testdata/no_permissios.txt", "/tmp/out.txt", 0, 0)
			require.Error(t, err)

			os.Chmod("testdata/no_permissios.txt", 0o420)
		})

		t.Run("no permission to write destination", func(t *testing.T) {
			err := Copy("testdata/input.txt", "/no_permissios.txt", 0, 0)
			require.Error(t, err)
		})

		t.Run("source is directory", func(t *testing.T) {
			err := Copy("testdata", "/tmp/out.txt", 0, 0)
			require.Error(t, err)
		})
	})

	t.Run("equal files", func(t *testing.T) {
		f1, err := os.Open("testdata/input.txt")
		require.NoError(t, err)
		defer f1.Close()

		s1, err := f1.Stat()
		require.NoError(t, err)

		f2, err := os.Open("testdata/input.txt")
		require.NoError(t, err)

		s2, err := f2.Stat()
		require.NoError(t, err)

		require.True(t, os.SameFile(s1, s2))
	})
}
