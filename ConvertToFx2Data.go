package main

import (
	"encoding/json"
	"strconv"
	"strings"
)

func preConvert(fx FxStruct, sigpathElement *SigpathElement) (err error) {

	err = nil
	id := strings.ToLower(fx.Descriptor)
	switch id {
	case "biasamp":
		sigpathElement.AmpType = "AmpHead"
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
		err = SplitAmpHeadCab(fx, fx2Preset)
	default:

	}
	return
}

// ConvertToFx2Data : convert fx1 to fx2
func ConvertToFx2Data(fx1Preset FX1Preset, presetName string) (js []byte, err error) {

	// init fx2 preset
	var fx2Preset Fx2PresetData
	fx2Preset.Bpm = 120
	fx2Preset.Version = "1.8"
	fx2Preset.Scenes.LatestEditedSceneSlot = -1
	fx2Preset.Scenes.Slot = make([]interface{}, 0)
	fx2Preset.Name = presetName

	// sigpath
	for _, fx := range fx1Preset.Fxs.Fx {
		var sigpathElement SigpathElement

		// pre convert
		err = preConvert(fx, &sigpathElement)

		if err != nil {
			return
		}

		// convert start
		sigpathElement.ModulePresetName = ""
		sigpathElement.Active, err = strconv.ParseBool(fx.Active)

		if err != nil {
			return
		}

		sigpathElement.DspID = fx.Descriptor
		sigpathElement.DspID = strings.Replace(sigpathElement.DspID, "LIVE.", "FX2.", 1)
		sigpathElement.ID = fx.Uniqueid
		sigpathElement.Selected, err = strconv.ParseBool(fx.Selected)

		// omitempty
		sigpathElement.AmpID = fx.AmpID
		sigpathElement.DistortionID = fx.DistortionID
		sigpathElement.DelayID = fx.DelayID
		sigpathElement.DodID = fx.DodID

		if err != nil {
			return
		}

		// copy parameters
		for _, param := range fx.Parameters.Parameter {

			var fx2Param SigpathParam
			fx2Param.ID, err = strconv.Atoi(param.Index)

			if err != nil {
				return
			}

			fx2Param.Value, err = strconv.ParseFloat(param.Value, 64)

			if err != nil {
				return
			}

			sigpathElement.Param = append(sigpathElement.Param, fx2Param)
		}

		// convert to json and append to sigpath
		fx2Preset.SigPath = append(fx2Preset.SigPath, sigpathElement)

		// post convert
		err = postConvert(fx, &fx2Preset)
		if err != nil {
			return
		}
	}

	js, _ = json.MarshalIndent(fx2Preset, "", "    ")
	return
}
