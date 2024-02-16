package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	builder := strings.Builder{}
	var err error = nil
	var idx = 0
	for i, v := range str {
		//Проверяем является ли элемент числом
		count, e := strconv.Atoi(string(str[i]))
		if e == nil && i == 0 { // Число стоит первым? Ошибка!
			err = ErrInvalidString
			break
		} else if e == nil && i > 0 { //Число не стоит первым? ок
			_, e2 := strconv.Atoi(string(str[i-1])) //А что стоит перед ним?
			if e == nil && e2 == nil {              // Два числа подряд? Ошибка!
				err = ErrInvalidString
				break
			}
		}
		if e == nil {
			if count != 0 {
				builder.WriteString(strings.Repeat(string(str[i-1]), count-1))
			}
			if count == 0 {
				idx++ // Фиксируем отклонение индекса
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
