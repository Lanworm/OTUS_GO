package hw09structvalidator

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	reIn    = regexp.MustCompile(`(in:|,)`)
	reRe    = regexp.MustCompile(`[^regexp:].*`)
	reDigit = regexp.MustCompile(`([0-9-]+)`)
)

type FieldValidator struct {
	field  reflect.StructField
	value  reflect.Value
	errors ValidationErrors
}

func (fv FieldValidator) validate() ValidationErrors {
	if keys, ok := fv.field.Tag.Lookup("validate"); ok && keys != "" {
		for _, key := range strings.Split(keys, "|") {
			var err error

			switch {
			case strings.Contains(key, "len:"):
				err = fv.validateLen(key)
			case strings.Contains(key, "regexp:"):
				err = fv.validateByRegexp(key)
			case strings.Contains(key, "min:"):
				err = fv.validateMin(key)
			case strings.Contains(key, "max:"):
				err = fv.validateMax(key)
			case strings.Contains(key, "in:"):
				err = fv.validateKeyIn(key)
			}

			if err == nil {
				continue
			}

			fv.errors = append(fv.errors, ValidationError{
				Field: fv.field.Name,
				Err:   err,
			})
		}
	}

	return fv.errors
}

func (fv FieldValidator) getSlice() (reflect.Value, error) {
	slice := reflect.ValueOf([]string{})

	switch fv.value.Kind() { //nolint:exhaustive
	case reflect.String:
		slice = reflect.Append(slice, reflect.ValueOf(fv.value.String()))
	case reflect.Int:
		slice = reflect.ValueOf([]int{})
		slice = reflect.Append(slice, fv.value)
	case reflect.Slice:
		slice = reflect.AppendSlice(slice, fv.value.Slice(0, fv.value.Len()))
	default:
		return slice, ErrValidateBaseRule
	}

	return slice, nil
}

func (fv FieldValidator) validateLen(key string) error {
	valKey, err := strconv.Atoi(reDigit.FindString(key))
	if err != nil {
		return err
	}

	slice, err := fv.getSlice()
	if err != nil {
		return err
	}

	for i := 0; i < slice.Len(); i++ {
		v := slice.Index(i)

		switch {
		case valKey < 0:
			return ErrValidateBaseRule
		case len(v.String()) != valKey:
			return ErrValidate
		}
	}

	return nil
}

func (fv FieldValidator) validateKeyIn(key string) error {
	arr := reIn.Split(key, -1)
	if len(arr) < 2 {
		return ErrValidateBaseRule
	}
	arr = arr[1:]

	slice, err := fv.getSlice()
	if err != nil {
		return err
	}

	var valResult string

	for i := 0; i < slice.Len(); i++ {
		val := slice.Index(i)

		switch val.Kind() { //nolint:exhaustive
		case reflect.String:
			valResult = val.String()
		case reflect.Int:
			valResult = strconv.Itoa(int(val.Int()))
		default:
			return ErrValidateBaseRule
		}

		if err := checkValidArrStr(valResult, arr); err != nil {
			return err
		}
	}

	return nil
}

func checkValidArrStr(value string, arr []string) error {
	for _, k := range arr {
		if k == value {
			return nil
		}
	}

	return ErrValidate
}

func (fv FieldValidator) validateByRegexp(key string) error {
	if fv.value.Kind() != reflect.String {
		return nil
	}

	re := reRe.FindString(key)
	result, err := regexp.MatchString(re, fv.value.String())
	if err != nil {
		return err
	}

	if !result {
		return ErrValidate
	}

	return nil
}

func (fv FieldValidator) validateInt(key string) (int, error) {
	if fv.value.Kind() != reflect.Int {
		return 0, ErrValidateBaseRule
	}

	valKey, err := strconv.Atoi(reDigit.FindString(key))
	if err != nil {
		return 0, err
	}

	if valKey < 0 {
		return 0, ErrValidateBaseRule
	}

	return valKey, nil
}

func (fv FieldValidator) validateMin(key string) error {
	value, err := fv.validateInt(key)
	if err != nil {
		return err
	}

	if fv.value.Int() < int64(value) {
		return ErrValidate
	}

	return nil
}

func (fv FieldValidator) validateMax(key string) error {
	value, err := fv.validateInt(key)
	if err != nil {
		return err
	}

	if fv.value.Int() > int64(value) {
		return ErrValidate
	}

	return nil
}
