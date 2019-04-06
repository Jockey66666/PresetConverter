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

// Fx2Struct : internal fx
type Fx2Struct struct {
	ModulePresetName string         `json:"ModulePresetName,omitempty"`
	Active           bool           `json:"active"`
	DspID            string         `json:"dspId"`
	ID               string         `json:"id,omitempty"`
	Selected         bool           `json:"selected"`
	AmpType          string         `json:"AmpType,omitempty"`
	AmpID            string         `json:"ampId,omitempty"`
	DistortionID     string         `json:"distortionId,omitempty"`
	DelayID          string         `json:"delayId,omitempty"`
	ModID            string         `json:"modId,omitempty"`
	Param            []SigpathParam `json:"param,omitempty"`
}

// InnerSigPathElement : sigpaths
type InnerSigPathElement struct {
	Fx      []Fx2Struct `json:"Fx"`
	Enabled bool        `json:"enabled"`
	Index   int         `json:"index"`
}

// DualPathElement : for dualpath
type DualPathElement struct {
	SigPaths []InnerSigPathElement `json:"sigPaths"`
	Type     string                `json:"type"`

	//middle fx
	Fx []Fx2Struct `json:"Fx,omitempty"`
}

// MixerChannel : mixer channels
type MixerChannel struct {
	Index int     `json:"index"`
	Delay float64 `json:"delay"`
	Gain  float64 `json:"gain"`
	Pan   float64 `json:"pan"`
}

// SigpathElement : sigpath data
type SigpathElement struct {
	ModulePresetName string         `json:"ModulePresetName,omitempty"`
	Active           bool           `json:"active"`
	DspID            string         `json:"dspId"`
	ID               string         `json:"id,omitempty"`
	Selected         bool           `json:"selected"`
	AmpType          string         `json:"AmpType,omitempty"`
	AmpID            string         `json:"ampId,omitempty"`
	DistortionID     string         `json:"distortionId,omitempty"`
	DelayID          string         `json:"delayId,omitempty"`
	ModID            string         `json:"modId,omitempty"`
	Param            []SigpathParam `json:"param,omitempty"`

	// for dual path
	SigPath []DualPathElement `json:"sigPath,omitempty"`

	// for mixer
	Channels []MixerChannel `json:"channels,omitempty"`
	Linked   bool           `json:"linked,omitempty"`
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
