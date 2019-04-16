package main

import (
	"encoding/json"
	"time"
)

// CreateMeta : create meta.json
func CreateMeta(path string, presetName string, author string, checker LicenseChecker) []byte {
	var metaJSON Fx2PresetMeta
	metaJSON.Author = author
	metaJSON.IsNew = true
	metaJSON.Name = presetName

	currentTime := time.Now()
	metaJSON.AddTime = currentTime.Format(time.RFC3339Nano)
	metaJSON.ModifiedTime = currentTime.Format(time.RFC3339Nano)

	metaJSON.LicenseInfo = checker.GetLicenseMeta(path + "/data.json")

	js, _ := json.MarshalIndent(metaJSON, "", "    ")
	return js
}
