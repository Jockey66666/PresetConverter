package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)

func createMeta(path string) []byte {
	var metaJSON Fx2PresetMeta
	metaJSON.Author = ""
	metaJSON.IsNew = true
	metaJSON.Name = "Sample"

	currentTime := time.Now()
	metaJSON.AddTime = currentTime.Format(time.RFC3339Nano)
	metaJSON.ModifiedTime = currentTime.Format(time.RFC3339Nano)

	metaJSON.LicenseInfo.Fx2LE = false
	metaJSON.LicenseInfo.Fx2License = 0
	metaJSON.LicenseInfo.ExpansionAcoustic = false
	metaJSON.LicenseInfo.ExpansionBass = false
	metaJSON.LicenseInfo.ExpansionMetal = false
	metaJSON.LicenseInfo.ModernVintage = false

	js, _ := json.MarshalIndent(metaJSON, "", "    ")
	return js
}

func copyThumbnail(path string) {

	input, err := ioutil.ReadFile("input/thumbnail.png")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = SaveFile(path+"/thumbnail.png", input)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// MigrationCore : mainly convert function
func MigrationCore(author string, presetSlice []PresetSliceStruct) int {
	count := 0
	for i, preset := range presetSlice {

		// step 1. read fx1 preset xml
		var data []byte
		var err error
		data, _ = OpenFile(preset.PresetPath)

		if err != nil {
			fmt.Printf("error: %v", err)
			continue
		}

		var fx1Preset FX1Preset
		err = xml.Unmarshal(data, &fx1Preset)

		if err != nil {
			fmt.Printf("error: %v", err)
			continue
		}

		// step 2. create preset folder
		uuid := uuid.Must(uuid.NewRandom())
		outputPath := "temp/" + preset.BankUUID + "/" + strings.ToUpper(uuid.String())
		CreateDirIfNotExist(outputPath)

		// step 3. create data.json
		dataJSON, err := ConvertToFx2Data(fx1Preset, preset.PresetName)

		if len(dataJSON) > 0 && err == nil {
			SaveFile(outputPath+"/data.json", dataJSON)
			count++
			presetSlice[i].MigrateResult = 1 // success
		} else {
			fmt.Println(err, "migrate", preset.PresetPath, "failed")
			presetSlice[i].MigrateResult = -1 // failed

			// remove folder and process next preset
			os.RemoveAll(outputPath)
			continue
		}

		// step 4. create meta.json
		metaJSON := createMeta(outputPath)
		SaveFile(outputPath+"/meta.json", metaJSON)
		// step 5. copy thumbnail.png
		copyThumbnail(outputPath)
	}
	return count
}
