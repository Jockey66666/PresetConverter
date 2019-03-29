package main

import (
	"fmt"
	"io/ioutil"
	"os"
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
