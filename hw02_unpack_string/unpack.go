package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	builder := strings.Builder{}

	var err error

	idx := 0

	for i, v := range str {
		count, e := strconv.Atoi(string(str[i]))

		if e == nil && i == 0 {
			err = ErrInvalidString
			break
		} else if e == nil && i > 0 {
			_, e2 := strconv.Atoi(string(str[i-1]))
			if e == nil && e2 == nil {
				err = ErrInvalidString
				break
			}
		}

		if e == nil {
			if count != 0 {
				builder.WriteString(strings.Repeat(string(str[i-1]), count-1))
			}
			if count == 0 {
				idx++
				result := builder.String()
				result = result[:i-idx]
				builder.Reset()
				builder.WriteString(result)
			}
		} else {
			builder.WriteString(string(v))
		}
	}

	return builder.String(), err
}
