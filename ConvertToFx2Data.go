package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func initFx2Preset(presetName string) Fx2PresetData {
	var fx2Preset Fx2PresetData
	fx2Preset.Bpm = 120
	fx2Preset.Version = "1.8"
	fx2Preset.Scenes.LatestEditedSceneSlot = -1
	fx2Preset.Scenes.Slot = make([]interface{}, 0)
	fx2Preset.Name = presetName

	return fx2Preset
}

// ConvertToFx2Data : convert fx1 to fx2
func ConvertToFx2Data(fx1Preset FX1Preset, presetName string) (js []byte, err error) {

	// init fx2 preset
	fx2Preset := initFx2Preset(presetName)

	/*
		only one group splitter and mixer for dual path
	*/

	dualPath := make(map[string]FxElement)

	for _, fx := range fx1Preset.Fxs.FxElements {

		elementName := strings.ToLower(fx.XMLName.Local)

		switch elementName {
		case "fx":
			ConvertFx(fx, &fx2Preset.SigPath, &fx2Preset.Embedded)
		case "live.splitter", "splitter":
			dualPath["splitter"] = fx
		case "live.mixer", "mixer":

			// check splitter is exist
			_, exists := dualPath["splitter"]
			if !exists {
				fmt.Println(presetName, "miss splitter?")
				return
			}

			dualPath["mixer"] = fx
			ConvertDualPath(dualPath, &fx2Preset.SigPath, &fx2Preset.Embedded)
		default:
			fmt.Println("unknown fx element.")
		}
	}

	js, _ = json.MarshalIndent(fx2Preset, "", "    ")
	return
}
