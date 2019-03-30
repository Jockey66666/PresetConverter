package main

// Fx2PresetData :  data.json
type Fx2PresetData struct {
	Bpm         int    `json:"bpm"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Name        string `json:"name"`
	Scenes      struct {
		LatestEditedSceneSlot int           `json:"latestEditedSceneSlot"`
		Slot                  []interface{} `json:"slot"`
	} `json:"scenes"`
	SigPath  []SigpathElement  `json:"sigPath"`
	Embedded []EmbeddedAmpData `json:"embedded"`
}

// SigpathParam :  parameters
type SigpathParam struct {
	ID    int     `json:"id"`
	Value float64 `json:"value"`
}

// SigpathElement : sigpath data
type SigpathElement struct {
	AmpType          string         `json:"AmpType,omitempty"`
	ModulePresetName string         `json:"ModulePresetName"`
	Active           bool           `json:"active"`
	AmpID            string         `json:"ampId,omitempty"`
	DspID            string         `json:"dspId"`
	ID               string         `json:"id"`
	Param            []SigpathParam `json:"param"`
	Selected         bool           `json:"selected"`
	DistortionID     string         `json:"distortionId,omitempty"`
	DelayID          string         `json:"delayId,omitempty"`
	DodID            string         `json:"modId,omitempty"`
}

// EmbeddedAmpData : ampdata
type EmbeddedAmpData struct {
	EmbeddedType string      `json:"EmbeddedType"`
	AmpID        string      `json:"ampId"`
	ID           string      `json:"id"`
	AmpData      interface{} `json:"ampData"`
	MetaData     interface{} `json:"metaData"`
	PanelData    interface{} `json:"panelData"`
}
