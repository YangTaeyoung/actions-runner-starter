package validator

import (
	"errors"
)

func Positive(ans interface{}) error {
	if ans.(int) <= 0 {
		return errors.New("it must be positive")
	}

	return nil
}
