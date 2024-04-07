package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	src, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer src.Close()

	stat, err := src.Stat()
	if err != nil {
		return err
	}
	if offset > stat.Size() {
		return ErrOffsetExceedsFileSize
	}
	available := stat.Size()
	if offset > 0 {
		available -= offset
	}
	if limit > 0 {
		if limit < available {
			available = limit
		}
	}

	dst, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if offset > 0 {
		_, err := src.Seek(offset, 0)
		if err != nil {
			return err
		}
	}

	reader := io.LimitReader(src, available)
	bar := pb.Full.Start64(available)
	barReader := bar.NewProxyReader(reader)
	_, err = io.CopyN(dst, barReader, available)
	if err != nil {
		return err
	}
	bar.Finish()

	return nil
}
