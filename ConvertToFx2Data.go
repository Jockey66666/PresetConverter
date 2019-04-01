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

	// sigpath
	for _, fx := range fx1Preset.Fxs.FxElements {

		elementName := strings.ToLower(fx.XMLName.Local)

		switch elementName {
		case "fx":
			ConvertFx(fx, &fx2Preset)
		case "live.splitter", "splitter":
			fmt.Println(fx.Descriptor)
		case "live.mixer", "mixer":
			fmt.Println(fx.Descriptor)
		default:
			fmt.Println("unknown fx element.")
		}
	}

	js, _ = json.MarshalIndent(fx2Preset, "", "    ")
	return
}
