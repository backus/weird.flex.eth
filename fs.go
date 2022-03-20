package main

import (
	"os"
	"path"
)

type FileTypeCheck byte

const (
	IsFile FileTypeCheck = iota
	IsDir
	DoesNotExist
	Unknown // Returned if stat fails and we're actually serving an error
)

func checkPathType(path string) (FileTypeCheck, error) {
	fileInfo, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			return DoesNotExist, nil
		} else {
			return Unknown, err
		}
	}

	if fileInfo.IsDir() {
		return IsDir, nil
	} else {
		return IsFile, nil
	}
}

func EnsureDirExists(path string) error {
	pathType, err := checkPathType(path)

	if err != nil {
		return err
	}

	if pathType == DoesNotExist {
		return os.Mkdir(path, 0777)
	} else {
		return nil
	}
}

func JoinProjectPath(relativePath string) (string, error) {
	projectDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return path.Join(projectDir, relativePath), nil
}
