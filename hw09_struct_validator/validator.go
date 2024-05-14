package hw09structvalidator

import (
	"reflect"
)

func Validate(v interface{}) error {
	errValid := ValidationErrors{}
	valStruct := reflect.ValueOf(v)

	if valStruct.Kind() != reflect.Struct {
		return ErrNotStructProvided
	}

	typeStruct := valStruct.Type()
	for i := 0; i < typeStruct.NumField(); i++ {
		fv := FieldValidator{
			field: typeStruct.Field(i),
			value: valStruct.Field(i),
		}

		errors := fv.validate()
		errValid = append(errValid, errors...)
	}

	if len(errValid) != 0 {
		return errValid
	}

	return nil
}
