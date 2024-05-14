package hw09structvalidator

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrNotStructProvided = errors.New("only struct provided")
	ErrValidateBaseRule  = errors.New("the data violates the underlying conditions")
	ErrValidate          = errors.New("validate error")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	strResult := strings.Builder{}

	for _, err := range v {
		if !errors.Is(err.Err, ErrValidate) {
			continue
		}

		strResult.WriteString(fmt.Sprintf("FieldValidator: %s: %s\n", err.Field, err.Err))
	}

	return strResult.String()
}
