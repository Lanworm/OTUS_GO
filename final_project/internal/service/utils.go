package service

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/Lanworm/OTUS_GO/final_project/internal/validation"
)

var (
	ErrInvalidNumberOfArguments           = errors.New("wrong number of arguments")
	ErrInvalidFormatOfArguments           = errors.New("wrong format of arguments")
	ErrInvalidArgumentTypeOfWidthOrHeight = errors.New("invalid argument type of width or height")
	ErrInvalidURL                         = errors.New("invalid URL")
)

func PrepareImgParams(u *url.URL) (imgParams *ImgParams, err error) {
	v := strings.Split(u.String(), "/")

	if len(v) < 4 {
		return nil, ErrInvalidNumberOfArguments
	}

	width := v[2]
	height := v[3]
	urlStr := strings.Join(v[4:], "/")

	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "http://" + urlStr
	}

	_, err = url.ParseRequestURI(urlStr)
	if err != nil {
		return nil, ErrInvalidURL
	}

	params, err := NewImgParams(width, height, urlStr)
	if err != nil {
		return nil, ErrInvalidFormatOfArguments
	}
	return params, nil
}

func NewImgParams(width string, height string, url string) (*ImgParams, error) {
	w, errw := strconv.Atoi(width)
	h, errh := strconv.Atoi(height)

	if errw != nil || errh != nil {
		return nil, ErrInvalidArgumentTypeOfWidthOrHeight
	}

	p := &ImgParams{
		Width:  uint(w),
		Height: uint(h),
		URL:    url,
	}

	err := validation.Validate(p)
	if err != nil {
		return nil, err
	}
	return p, nil
}
