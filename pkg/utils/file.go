package utils

import (
	"os"
	"path"
)

func ReadFileFrom(filepath string, rootpath ...string) (string, error) {
	var err error
	basepath := "./"
	if len(rootpath) > 0 && rootpath[0] != "" {
		basepath = rootpath[0]
	} else {
		basepath, err = os.Getwd()
		if err != nil {
			return "", WrapError(err, "failed to get working directory")
		}
	}
	fullpath := path.Join(basepath, filepath)
	data, err := os.ReadFile(fullpath)
	if err != nil {
		return "", WrapError(err, "failed to read file from %s", fullpath)
	}
	return string(data), nil
}

func OpenFileFrom(filepath string, rootpath ...string) (*os.File, error) {
	var err error
	basepath := "./"
	if len(rootpath) > 0 && rootpath[0] != "" {
		basepath = rootpath[0]
	} else {
		basepath, err = os.Getwd()
		if err != nil {
			return nil, WrapError(err, "failed to get working directory")
		}
	}
	fullpath := path.Join(basepath, filepath)
	file, err := os.OpenFile(fullpath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, WrapError(err, "failed to open file from %s", fullpath)
	}
	return file, nil
}
