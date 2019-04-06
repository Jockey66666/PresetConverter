package main

// FxStruct : fx and amp
type FxStruct struct {
	BmMask               int      `json:"BmMask,omitempty"`
	CableStyle           string   `json:"CableStyle,omitempty"`
	Category             string   `json:"Category,omitempty"`
	Category2            string   `json:"category,omitempty"`
	DisplayName          string   `json:"DisplayName,omitempty"`
	DisplayName2         string   `json:"displayName,omitempty"`
	DisplayUnit          int      `json:"DisplayUnit,omitempty"`
	EnableForLiteEdition bool     `json:"EnableForLiteEdition"`
	PackID               []string `json:"PackId,omitempty"`
	PackID2              []string `json:"packId,omitempty"`
	SigPathHeight        float64  `json:"SigPathHeight,omitempty"`
	SigpathDescriptor    string   `json:"SigpathDescriptor"`
	StoreCategory        string   `json:"StoreCategory,omitempty"`
	StoreSortKey         string   `json:"StoreSortKey,omitempty"`
	Description          string   `json:"description,omitempty"`
	DspID                string   `json:"dspId"`
	IsFavorite           bool     `json:"isFavorite"`
	Hidden               bool     `json:"Hidden,omitempty"`
	Descriptor           string   `json:"descriptor,omitempty"`
	IsPresetFormat       bool     `json:"isPresetFormat,omitempty"`
	Subtitle             string   `json:"subtitle,omitempty"`
	LicenseTierReq       *int     `json:"LicenseTierReq,omitempty"`
}

// BiasFxJSON : bias fx 2 json file
type BiasFxJSON struct {
	BiasFx2 []struct {
		Category string     `json:"category"`
		FxList   []FxStruct `json:"fxList,omitempty"`
		AmpList  []FxStruct `json:"ampList,omitempty"`
	} `json:"BIAS_FX_2"`
	BiasAMP2 []struct {
		Category string     `json:"category"`
		FxList   []FxStruct `json:"fxList,omitempty"`
		AmpList  []FxStruct `json:"ampList,omitempty"`
	} `json:"BIAS_AMP_2,omitempty"`
	Version string `json:"version"`
}
