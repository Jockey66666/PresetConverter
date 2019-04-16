package main

import (
	"os"
	"runtime"
)

// PathGenerator : get migration tool paths
type PathGenerator struct {
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// DocPath : document path
func (gen *PathGenerator) DocPath() string {
	return userHomeDir() + "/Documents/PositiveGrid/BIAS_FX2/MigrationTool"
}

// InputSettingPath : Input setting path
func (gen *PathGenerator) InputSettingPath() string {
	return gen.DocPath() + "/input_bank_list.json"
}

// LocalSettingPath : local input setting path
func (gen *PathGenerator) LocalSettingPath() string {
	return "input_bank_list.json"
}
