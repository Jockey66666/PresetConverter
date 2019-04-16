package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/uuid"
)

func copyThumbnail(path string) {

	input, err := ioutil.ReadFile("thumbnail.png")
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
func MigrationCore(checker LicenseChecker, inputBanks InputBankListJSON, presetSlice []PresetSliceStruct) int {
	count := 0
	for i, preset := range presetSlice {

		// step 1. read fx1 preset xml
		var data []byte
		var err error
		data, _ = OpenFile(preset.PresetPath)

		if err != nil {
			fmt.Printf("error: %v", err)
			presetSlice[i].MigrateResult = -1 // failed
			continue
		}

		var fx1Preset FX1Preset
		err = xml.Unmarshal(data, &fx1Preset)

		if err != nil {
			fmt.Printf("error: %v", err)
			presetSlice[i].MigrateResult = -1 // failed
			continue
		}

		// step 2. create preset folder
		uuid := uuid.Must(uuid.NewRandom())
		outputPath := inputBanks.Temp + "/" + preset.BankUUID + "/" + strings.ToUpper(uuid.String())
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
		metaJSON := CreateMeta(outputPath, preset.PresetName, inputBanks.Author, checker)
		SaveFile(outputPath+"/meta.json", metaJSON)
		// step 5. copy thumbnail.png
		copyThumbnail(outputPath)
	}
	return count
}
