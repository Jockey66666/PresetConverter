package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/thedevsaddam/gojsonq"
	"github.com/tidwall/gjson"
)

//var globalPresetName = "Change amp"
var globalPresetName = "Sample"

func splictAmpHeadCab(fx FxStruct, fx2Preset *Fx2PresetData) (err error) {

	if len(fx.Ampdata) > 0 {
		//embedded data
		var ampData EmbeddedAmpData
		ampData.ID = fx.Uniqueid
		ampData.AmpID = fx.AmpID
		ampData.EmbeddedType = "BiasAmp"

		var extraData interface{}
		err = json.Unmarshal([]byte(fx.Ampdata), &extraData)

		if err != nil {
			fmt.Println(err)
			return
		}

		m := extraData.(map[string]interface{})

		ampData.AmpData = m["ampData"]
		ampData.MetaData = m["metaData"]
		ampData.PanelData = m["panelData"]

		fx2Preset.Embedded = append(fx2Preset.Embedded, ampData)

		// cab sigpath
		var sigpathElement SigpathElement
		sigpathElement.ModulePresetName = ""
		sigpathElement.Active, err = strconv.ParseBool(fx.Active)

		if err != nil {
			return
		}

		sigpathElement.AmpType = "AmpCab"
		sigpathElement.AmpID = fx.AmpID
		sigpathElement.DspID = "BiasAmp"
		sigpathElement.ID = fx.Uniqueid
		sigpathElement.Selected, err = strconv.ParseBool(fx.Selected)

		if err != nil {
			return
		}

		ampDataByte, _ := json.Marshal(ampData.AmpData)
		query := gojsonq.New().JSONString(string(ampDataByte)).
			From("sigPath.blocks.items").Where("id", "=", "bias.cab").Only("params")
		var queryString []byte
		queryString, err = json.Marshal(query)

		if err != nil {
			return
		}

		result := gjson.Get(string(queryString), "0.params")
		result.ForEach(func(key, value gjson.Result) bool {
			r1 := gjson.Parse(value.String()).Get("id")
			r2 := gjson.Parse(value.String()).Get("value")

			var fx2Param SigpathParam
			fx2Param.ID = int(r1.Int()) // int64 to int
			fx2Param.Value = r2.Float()
			sigpathElement.Param = append(sigpathElement.Param, fx2Param)

			return true
		})

		fx2Preset.SigPath = append(fx2Preset.SigPath, sigpathElement)
	}

	return
}

func preConvert(fx FxStruct, sigpathElement *SigpathElement) (err error) {

	err = nil
	id := strings.ToLower(fx.Descriptor)
	switch id {
	case "biasamp":
		sigpathElement.AmpType = "AmpHead"
		sigpathElement.AmpID = fx.AmpID
	default:
		// do nothing
	}
	return
}

func postConvert(fx FxStruct, fx2Preset *Fx2PresetData) (err error) {

	err = nil
	id := strings.ToLower(fx.Descriptor)
	switch id {
	case "biasamp":
		err = splictAmpHeadCab(fx, fx2Preset)
	default:

	}
	return
}

func convertToFx2Data(fx1Preset FX1Preset, fx2Preset Fx2PresetData) []byte {

	// sigpath
	for _, fx := range fx1Preset.Fxs.Fx {
		var sigpathElement SigpathElement
		var err error

		sigpathElement.ModulePresetName = ""
		sigpathElement.Active, err = strconv.ParseBool(fx.Active)

		err = preConvert(fx, &sigpathElement)

		if err != nil {
			continue
		}

		sigpathElement.DspID = fx.Descriptor
		sigpathElement.DspID = strings.Replace(sigpathElement.DspID, "LIVE.", "FX2.", 1)
		sigpathElement.ID = fx.Uniqueid
		sigpathElement.Selected, err = strconv.ParseBool(fx.Selected)

		if err != nil {
			continue
		}

		for _, param := range fx.Parameters.Parameter {

			var fx2Param SigpathParam
			fx2Param.ID, err = strconv.Atoi(param.Index)

			if err != nil {
				continue
			}

			fx2Param.Value, err = strconv.ParseFloat(param.Value, 64)

			if err != nil {
				continue
			}

			sigpathElement.Param = append(sigpathElement.Param, fx2Param)
		}

		// convert to json and append to sigpath
		fx2Preset.SigPath = append(fx2Preset.SigPath, sigpathElement)

		err = postConvert(fx, &fx2Preset)
		if err != nil {
			continue
		}
	}

	js, _ := json.MarshalIndent(fx2Preset, "", "    ")
	return js
}

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
		dataJSON := convertToFx2Data(fx1Preset, fx2PresetData)

		if len(dataJSON) > 0 {
			SaveFile(outputPath+"/data.json", dataJSON)
			count++
		}

		// meta.json
		metaJSON := createMeta(outputPath)
		SaveFile(outputPath+"/meta.json", metaJSON)

		// thumbnail.png
		copyThumbnail(outputPath)

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
