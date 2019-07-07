package validator

import "errors"

var maxMsgLength = 10

func Length(msg string) (bool, error) {
	validation := len(msg) < maxMsgLength
	if !validation {
		return false, errors.New("Mensaje muy largo")
	}
	return true, nil
}

func MinLength(msg string, min int) (bool, error) {
	validation := len(msg) < min
	if validation {
		return false, errors.New("Mensaje muy corto")
	}
	return true, nil
}

func LengthOfParameters(params []string) (bool, error) {
	for _, s := range params {
		v, err := Length(s)
		if !v {
			return false, err
		}
	}
	return true, nil
}
