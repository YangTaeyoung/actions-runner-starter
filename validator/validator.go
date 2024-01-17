package validator

import (
	"errors"
	"strconv"
)

func Positive(ans interface{}) error {
	ansStr, ok := ans.(string)
	if !ok {
		return errors.New("it must be string")
	}

	ansInt, err := strconv.Atoi(ansStr)
	if err != nil {
		return err
	}

	if ansInt <= 0 {
		return errors.New("it must be positive")
	}

	return nil
}
