package main

/*

```

├── GlobalPresets
|   ├── bank.json
|   └── (banks)
└── MigrationTool
    ├── report
    |    └── MigrationReport_MMDD_N.json
    ├── temp
    |    ├── bank.json (output)
    |    └── (banks) (output)
    |
    ├── input_bank_list.json (input)
    ├── PresetConverterWin.exe
    └── PresetConverterMac

```

*/

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"runtime"
	"time"
)

//var globalPresetName = "Change amp"
var globalPresetName = "Sample"
var globalDebug = false

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

func main() {
	t := time.Now()
	defer func() {
		elapsed := time.Since(t)
		fmt.Println("")
		fmt.Println("elapsed:", elapsed)
	}()

	// pre migration
	inputBanks, presetSlice, bankTable := PreMigration()

	// migration
	cpu := runtime.NumCPU()
	num := len(presetSlice)

	countChannel := make(chan int)
	for i := 0; i < cpu; i++ {
		from := i * num / cpu
		to := (i + 1) * num / cpu
		go func() {
			countChannel <- MigrationCore(inputBanks.Author, presetSlice[from:to])
		}()
	}

	total := 0
	for i := 0; i < cpu; i++ {
		total += <-countChannel
	}

	fmt.Println(total, "presets migrated")

	// post migration
	PostMigration(inputBanks.Dst, bankTable, presetSlice)

	/*
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

	*/
}
