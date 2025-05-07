package gemini

import (
	"os"
	"path"
)

func readPromptFile(filepath string, rootpath ...string) (string, error) {
	var err error
	basepath := "./"
	if len(rootpath) > 0 && rootpath[0] != "" {
		basepath = rootpath[0]
	} else {
		basepath, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}
	fullpath := path.Join(basepath, filepath)
	data, err := os.ReadFile(fullpath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func openAttachmentFile(filepath string, rootpath ...string) (*os.File, error) {
	var err error
	basepath := "./"
	if len(rootpath) > 0 && rootpath[0] != "" {
		basepath = rootpath[0]
	} else {
		basepath, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	fullpath := path.Join(basepath, filepath)
	return os.OpenFile(fullpath, os.O_RDONLY, 0644)
}
