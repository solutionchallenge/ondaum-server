package utils

import (
	"errors"
	"fmt"
)

func NewError(fmtstr string, args ...interface{}) error {
	return fmt.Errorf(fmtstr, args...)
}

func WrapError(err error, fmtstr string, args ...interface{}) error {
	return fmt.Errorf("%s: %w", fmt.Sprintf(fmtstr, args...), err)
}

func IsError(err error, target error) bool {
	return errors.Is(err, target)
}
