package main

import "strconv"

// ConvertDualPath : convert dual path
func ConvertDualPath(dualPath map[string]FxElement, sigpath *[]SigpathElement, embedded *[]EmbeddedData) {
	splitter := dualPath["splitter"]
	convertSplitter(splitter, sigpath, embedded)

	mixer := dualPath["mixer"]
	convertMixer(splitter, mixer, sigpath)
}

func convertSplitter(splitter FxElement, sigpath *[]SigpathElement, embedded *[]EmbeddedData) {
	sigpathElement := GetSigpathElement(splitter)
	sigpathElement.DspID = "FX2.Splitter"

	var dualPathElement DualPathElement
	for _, fxs := range splitter.Fxs {
		var innerSigpath InnerSigPathElement
		innerSigpath.Index, _ = strconv.Atoi(fxs.Index)
		innerSigpath.Enabled, _ = strconv.ParseBool(fxs.Enabled)

		for _, fx := range fxs.Fx {
			var tmp []SigpathElement
			ConvertFx(fx, &tmp, embedded)
			// copy to dualPath
			copyToDualSigpath(tmp, &innerSigpath)
		}

		dualPathElement.SigPaths = append(dualPathElement.SigPaths, innerSigpath)
	}

	dualPathElement.Type = "dualPath"
	sigpathElement.SigPath = append(sigpathElement.SigPath, dualPathElement)

	*sigpath = append(*sigpath, sigpathElement)
}

func copyToDualSigpath(sigpath []SigpathElement, innerSigpath *InnerSigPathElement) {
	for _, s := range sigpath {

		fx2 := Fx2Struct{
			ModulePresetName: s.ModulePresetName,
			Active:           s.Active,
			DspID:            s.DspID,
			ID:               s.ID,
			Selected:         s.Selected,
			AmpType:          s.AmpType,
			AmpID:            s.AmpID,
			DistortionID:     s.DistortionID,
			DelayID:          s.DelayID,
			ModID:            s.ModID,
			Param:            s.Param,
		}

		innerSigpath.Fx = append(innerSigpath.Fx, fx2)
	}
}

func convertMixer(splitter FxElement, mixer FxElement, sigpath *[]SigpathElement) {
	var sigpathElement SigpathElement
	sigpathElement.Active, _ = strconv.ParseBool(mixer.Active)
	sigpathElement.DspID = "FX2.Mixer"
	sigpathElement.Linked, _ = strconv.ParseBool(splitter.Linked)
	sigpathElement.Selected = false

	for _, fxs := range splitter.Fxs {
		var c MixerChannel
		c.Index, _ = strconv.Atoi(fxs.Index)
		c.Delay, _ = strconv.ParseFloat(fxs.Delay, 64)
		c.Gain, _ = strconv.ParseFloat(fxs.Gain, 64)
		c.Pan, _ = strconv.ParseFloat(fxs.Pan, 64)

		sigpathElement.Channels = append(sigpathElement.Channels, c)
	}
	*sigpath = append(*sigpath, sigpathElement)
}
