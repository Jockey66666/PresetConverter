package main

import (
	"strconv"
	"strings"
)

// GetSigpathElement : get single sigpath element from xml element
func GetSigpathElement(fx FxElement) (sigpathElement SigpathElement) {
	sigpathElement.ModulePresetName = ""
	sigpathElement.Active, _ = strconv.ParseBool(fx.Active)
	sigpathElement.DspID = fx.Descriptor
	sigpathElement.DspID = strings.Replace(sigpathElement.DspID, "LIVE.", "FX2.", 1)
	sigpathElement.ID = fx.Uniqueid
	sigpathElement.Selected, _ = strconv.ParseBool(fx.Selected)

	// copy parameters
	for _, param := range fx.Parameters.Parameter {
		var fx2Param SigpathParam
		fx2Param.ID, _ = strconv.Atoi(param.Index)
		fx2Param.Value, _ = strconv.ParseFloat(param.Value, 64)
		sigpathElement.Param = append(sigpathElement.Param, fx2Param)
	}

	return
}

// ConvertFx : for <Fx>
func ConvertFx(fx FxElement, fx2Preset *Fx2PresetData) {

	fxID := strings.ToLower(fx.Descriptor)
	switch fxID {
	case "biasamp":
		HandleAmp(fx, fx2Preset)
	default:
		sigpathElement := GetSigpathElement(fx)
		fx2Preset.SigPath = append(fx2Preset.SigPath, sigpathElement)
	}
}
