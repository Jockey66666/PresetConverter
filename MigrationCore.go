package main

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

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
			//SaveFile(outputPath+"/data.json", dataJSON)
			count++
			presetSlice[i].MigrateResult = 1 // success
		} else {
			fmt.Println(err, "migrate", preset.PresetPath, "failed")
			presetSlice[i].MigrateResult = -1 // failed
		}

		// step 4. create meta.json
		// step 5. copy thumbnail.png
	}
	return count
}
