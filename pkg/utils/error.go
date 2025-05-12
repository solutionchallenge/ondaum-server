package utils

import (
	"fmt"
)

func NewError(fmtstr string, args ...interface{}) error {
	_, filename, line := GetCallerInfo(2)
	return fmt.Errorf("(%s:%d) %w", filename, line, fmt.Errorf(fmtstr, args...))
}

func WrapError(err error, fmtstr string, args ...interface{}) error {
	_, filename, line := GetCallerInfo(2)
	return fmt.Errorf("(%s:%d) %s: %w", filename, line, fmt.Sprintf(fmtstr, args...), err)
}

func PassError(err error) error {
	_, filename, line := GetCallerInfo(2)
	return fmt.Errorf("(%s:%d): %w", filename, line, err)
}
