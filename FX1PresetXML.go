package main

import "encoding/xml"

// FX1Preset : preset xml
type FX1Preset struct {
	XMLName     xml.Name `xml:"iTonesPreset"`
	Text        string   `xml:",chardata"`
	Name        string   `xml:"name,attr"`
	Description string   `xml:"description,attr"`
	Category    string   `xml:"category,attr"`
	Version     string   `xml:"version,attr"`
	Fxs         struct {
		Text     string         `xml:",chardata"`
		Fx       []FxStruct     `xml:"Fx"`
		Splitter []LIVESplitter `xml:"LIVE.Splitter"`
		Mixer    []LIVEMixer    `xml:"LIVE.Mixer"`
	} `xml:"Fxs"`
}

// FxStruct : fx struct
type FxStruct struct {
	Text         string `xml:",chardata"`
	Active       string `xml:"active,attr"`
	Selected     string `xml:"selected,attr"`
	Uniqueid     string `xml:"uniqueid,attr"`
	Descriptor   string `xml:"descriptor,attr"`
	AmpID        string `xml:"ampId,attr"`
	DistortionID string `xml:"distortionId,attr"`
	DelayID      string `xml:"delayId,attr"`
	DodID        string `xml:"modId,attr"`
	Parameters   struct {
		Text      string `xml:",chardata"`
		Parameter []struct {
			Text  string `xml:",chardata"`
			Index string `xml:"index,attr"`
			Value string `xml:"value,attr"`
			Type  string `xml:"type,attr"`
		} `xml:"Parameter"`
	} `xml:"Parameters"`
	Ampdata string `xml:"Ampdata"`
}

// LIVESplitter : splitter
type LIVESplitter struct {
	Text       string `xml:",chardata"`
	Active     string `xml:"active,attr"`
	Selected   string `xml:"selected,attr"`
	Descriptor string `xml:"descriptor,attr"`
	Linked     string `xml:"linked,attr"`
	Parameters struct {
		Text      string `xml:",chardata"`
		Parameter []struct {
			Text  string `xml:",chardata"`
			Index string `xml:"index,attr"`
			Value string `xml:"value,attr"`
		} `xml:"Parameter"`
	} `xml:"Parameters"`
	Fxs []struct {
		Text    string     `xml:",chardata"`
		Index   string     `xml:"index,attr"`
		Pan     string     `xml:"pan,attr"`
		Gain    string     `xml:"gain,attr"`
		Delay   string     `xml:"delay,attr"`
		Enabled string     `xml:"enabled,attr"`
		Fx      []FxStruct `xml:"Fx"`
	} `xml:"Fxs"`
}

// LIVEMixer : mixer
type LIVEMixer struct {
	Text       string `xml:",chardata"`
	Active     string `xml:"active,attr"`
	Selected   string `xml:"selected,attr"`
	Descriptor string `xml:"descriptor,attr"`
}
