package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
	"time"

	"github.com/google/uuid"
)

//var globalPresetName = "Change amp"
var globalPresetName = "Sample"
var globalDebug = true

func createMeta(path string) []byte {
	var metaJSON Fx2PresetMeta
	metaJSON.Author = ""
	metaJSON.IsNew = true
	metaJSON.Name = globalPresetName

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

func processFunction(presetFileNames []string) int {
	count := 0
	for _, presetFileName := range presetFileNames {
		var fx2PresetData Fx2PresetData
		fx2PresetData.Bpm = 120
		fx2PresetData.Version = "1.8"
		fx2PresetData.Scenes.LatestEditedSceneSlot = -1
		fx2PresetData.Scenes.Slot = make([]interface{}, 0)
		fx2PresetData.Name = globalPresetName

		var data []byte
		var err error
		data, err = OpenFile(presetFileName)

		if err != nil {
			continue
		}

		var fx1Preset FX1Preset
		err = xml.Unmarshal(data, &fx1Preset)

		if err != nil {
			fmt.Printf("error: %v", err)
			continue
		}

		uuid := uuid.Must(uuid.NewRandom())

		outputPath := "output/" + strings.ToUpper(uuid.String())
		CreateDirIfNotExist(outputPath)

		// data.json
		dataJSON := ConvertToFx2Data(fx1Preset, fx2PresetData)

		if len(dataJSON) > 0 {
			SaveFile(outputPath+"/data.json", dataJSON)
			count++
		}

		if globalDebug == false {
			// meta.json
			metaJSON := createMeta(outputPath)
			SaveFile(outputPath+"/meta.json", metaJSON)

			// thumbnail.png
			copyThumbnail(outputPath)
		}

		//fmt.Println(string(dataJSON))
	}

	return count
}

func main() {
	t := time.Now()

	// input file

	presetFileNames := []string{}

	i := 0
	for i < 1 {
		presetFileNames = append(presetFileNames, "input/"+globalPresetName+".Preset")
		i++
	}

	countChannel := make(chan int)

	cpu := runtime.NumCPU()
	num := len(presetFileNames)

	for i := 0; i < cpu; i++ {
		from := i * num / cpu
		to := (i + 1) * num / cpu
		go func() {
			countChannel <- processFunction(presetFileNames[from:to])
		}()
	}

	total := 0
	for i := 0; i < cpu; i++ {
		total += <-countChannel
	}

	elapsed := time.Since(t)
	fmt.Println(total, "preset migrated.")
	fmt.Println("elapsed:", elapsed)
}
