package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// OpenFile : read file byte
func OpenFile(path string) ([]byte, error) {
	fileHandle, err := os.Open(path)

	if err != nil {
		fmt.Println(path, "open failed")
		return nil, err
	}

	defer fileHandle.Close()
	data, err := ioutil.ReadAll(fileHandle)

	if err != nil {
		fmt.Println(path + "read failed")
		return nil, err
	}

	return data, nil
}

// SaveFile : save file byte
func SaveFile(path string, data []byte) error {
	if globalDebug == true {
		fmt.Println(string(data))
		return nil
	}
	return ioutil.WriteFile(path, data, 0644)
}

// CreateDirIfNotExist : create directory
func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			panic(err)
		}
	}
}

// GetSubFolderList : get subfolder list, depth = 1
func GetSubFolderList(files *[]string) filepath.WalkFunc {

	foundRootDir := false
	currentDepth := -1

	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}

		if foundRootDir {

			if currentDepth < 0 {
				currentDepth = strings.Count(path, "/")
			}

			if currentDepth < strings.Count(path, "/") {
				return filepath.SkipDir
			}
		}

		if info.IsDir() {

			if foundRootDir {
				*files = append(*files, info.Name())
			} else {
				foundRootDir = true
			}

		}

		return nil
	}
}
