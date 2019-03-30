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
		err = SplitAmpHeadCab(fx, fx2Preset)
	default:

	}
	return
}

func ConvertToFx2Data(fx1Preset FX1Preset, fx2Preset Fx2PresetData) []byte {

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
