package main

import (
	"encoding/json"
	"fmt"
)

// HandleDistortion : handle bias pedal distortion data
func HandleDistortion(fx FxElement, sigpath *[]SigpathElement, embedded *[]EmbeddedData) {

	sigpathElement := GetSigpathElement(fx)
	sigpathElement.DistortionID = fx.DistortionID

	var data EmbeddedData
	data.EmbeddedType = "BiasPedal_Distortion"
	data.DistortionID = fx.DistortionID
	data.ID = fx.Uniqueid

	appendEmbeddedData(fx.Distortiondata, data, embedded)
	*sigpath = append(*sigpath, sigpathElement)
}

// HandleDelay : handle bias pedal distortion data
func HandleDelay(fx FxElement, sigpath *[]SigpathElement, embedded *[]EmbeddedData) {

	sigpathElement := GetSigpathElement(fx)
	sigpathElement.DelayID = fx.DelayID

	var data EmbeddedData
	data.EmbeddedType = "BiasPedal_Delay"
	data.DelayID = fx.DelayID
	data.ID = fx.Uniqueid

	appendEmbeddedData(fx.Delaydata, data, embedded)
	*sigpath = append(*sigpath, sigpathElement)
}

// HandleModulation : handle bias pedal distortion data
func HandleModulation(fx FxElement, sigpath *[]SigpathElement, embedded *[]EmbeddedData) {

	sigpathElement := GetSigpathElement(fx)
	sigpathElement.ModID = fx.ModID

	var data EmbeddedData
	data.EmbeddedType = "BiasPedal_Modulaton"
	data.ModID = fx.ModID
	data.ID = fx.Uniqueid

	appendEmbeddedData(fx.Moddata, data, embedded)
	*sigpath = append(*sigpath, sigpathElement)
}

func appendEmbeddedData(fx1Data string, fx2Data EmbeddedData, embedded *[]EmbeddedData) {

	if len(fx1Data) > 0 {
		var extraData interface{}
		err := json.Unmarshal([]byte(fx1Data), &extraData)

		if err != nil {
			fmt.Println(err)
			return
		}

		m := extraData.(map[string]interface{})

		fx2Data.DistortionData = m["distortionData"]
		fx2Data.DelayData = m["delayData"]
		fx2Data.ModData = m["modData"]
		fx2Data.MetaData = m["metaData"]
		fx2Data.PanelData = m["panelData"]

		*embedded = append(*embedded, fx2Data)
	}
}
