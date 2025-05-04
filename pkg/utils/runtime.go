package utils

import (
	"runtime"
	"strings"
)

func GetCallerInfo(skip int) (string, string, int) {
	_, file, line, ok := runtime.Caller(skip + 1)
	if !ok {
		return "unknown", "unknown", 0
	}
	parts := strings.Split(file, "/")
	path := strings.Join(parts[:len(parts)-1], "/")
	filename := parts[len(parts)-1]
	return path, filename, line
}
