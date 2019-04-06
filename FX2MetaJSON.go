package main

// Fx2PresetMeta : meta.json
type Fx2PresetMeta struct {
	AddTime      string         `json:"add_time"`
	Author       string         `json:"author"`
	IsNew        bool           `json:"is_new"`
	LicenseInfo  LicenseTierReq `json:"licenseTierReq"`
	ModifiedTime string         `json:"modified_time"`
	Name         string         `json:"name"`
}

// LicenseTierReq : license info
type LicenseTierReq struct {
	Amp2LE            *bool `json:"Amp2LE,omitempty"`
	Amp2License       *int  `json:"Amp2License,omitempty"`
	Fx2LE             bool  `json:"Fx2LE"`
	Fx2License        int   `json:"Fx2License"`
	ExpansionAcoustic bool  `json:"expansion_acoustic"`
	ExpansionBass     bool  `json:"expansion_bass"`
	ExpansionMetal    bool  `json:"expansion_metal"`
	ModernVintage     bool  `json:"modern_vintage"`
}
