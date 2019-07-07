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

func MaxLength(msg string, max int) error {
	validation := len(msg) > max
	if validation {
		return errors.New("Mensaje muy largo")
	}
	return nil
}

func MinLength(msg string, min int) error {
	validation := len(msg) < min
	if validation {
		return errors.New("Mensaje muy corto")
	}
	return nil
}

func CheckErrors(errors []error) error {
	for _, err := range errors{
		if err != nil {
			return err
		}
	}
	return nil
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
