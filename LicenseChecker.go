package main

import (
	"encoding/json"
)

const (
	metalPack    = 0
	acousticPack = 1
	bassPack     = 2
)

var licenseTable = map[string]int{
	"Demo":  0,
	"Lite":  1,
	"Std":   2,
	"Pro":   3,
	"ELite": 4,
}

var extensionPack = map[string]int{
	"MTS9EAD6DBA54EDAADA20D97EF2FAAA1": metalPack,
	"ACSFD9C18B744E7080AFAABE7049FB4B": acousticPack,
	"BAS17A9D0E2841E0915D0E4FD63C437F": bassPack,
}

// LicenseChecker : generate LicenseTierReq for meta.json
type LicenseChecker struct {
	fx2Licenses  map[string]int
	amp2Licenses map[string]int
	packLicenses map[string]int
}

// Init : init license tier map
func (checker *LicenseChecker) Init(path string) (err error) {
	checker.fx2Licenses = make(map[string]int)
	checker.packLicenses = make(map[string]int)
	checker.amp2Licenses = make(map[string]int)

	err = initFx2Map(path+"fxlist.json", checker.fx2Licenses, checker.packLicenses)

	if err != nil {
		return
	}

	err = initAmp2Map(path+"amplist.json", checker.fx2Licenses, checker.amp2Licenses, checker.packLicenses)
	return
}

// GetLicenseMeta : get LicenseTierReq from data.json
func (checker LicenseChecker) GetLicenseMeta(presetPath string) (licenseTierReq LicenseTierReq) {
	data, _ := OpenFile(presetPath)
	var presetData Fx2PresetData
	err := json.Unmarshal(data, &presetData)
	if err != nil {
		return
	}

	// get dspid and ampid
	sigpathIDs := getSigpathIDs(presetData)

	// compare with license maps
	for _, id := range sigpathIDs {

		// find amp2
		if value, found := checker.amp2Licenses[id]; found {

			if licenseTierReq.Amp2LE == nil {
				licenseTierReq.Amp2LE = new(bool)
			}

			*licenseTierReq.Amp2LE = false

			if licenseTierReq.Amp2License == nil {
				licenseTierReq.Amp2License = new(int)
			}

			if *licenseTierReq.Amp2License < value {
				*licenseTierReq.Amp2License = value
			}
		}

		// find fx2
		if value, found := checker.fx2Licenses[id]; found {
			if licenseTierReq.Fx2License < value {
				licenseTierReq.Fx2License = value
			}
		}

		// find extension pack
		if value, found := checker.packLicenses[id]; found {
			switch value {
			case metalPack:
				licenseTierReq.ExpansionMetal = true
			case acousticPack:
				licenseTierReq.ExpansionAcoustic = true
			case bassPack:
				licenseTierReq.ExpansionBass = true
			}
		}
	}

	return
}

func getSigpathIDs(presetData Fx2PresetData) (sigpathIDs []string) {
	for _, element := range presetData.SigPath {
		//dual path
		if len(element.SigPath) > 0 {
			for _, inner := range element.SigPath {
				for _, p := range inner.SigPaths {
					for _, fx := range p.Fx {
						collectIDs(fx, &sigpathIDs)
					}
				}

				// middle fx
				for _, fx := range inner.Fx {
					collectIDs(fx, &sigpathIDs)
				}
			}
		} else {
			tmp := Fx2Struct{
				AmpID:        element.AmpID,
				DistortionID: element.DistortionID,
				DelayID:      element.DelayID,
				ModID:        element.ModID,
				DspID:        element.DspID,
			}
			collectIDs(tmp, &sigpathIDs)
		}
	}
	return
}

func collectIDs(fx Fx2Struct, sigpathIDs *[]string) {
	if len(fx.AmpID) > 0 {
		*sigpathIDs = append(*sigpathIDs, fx.AmpID)
	} else if len(fx.DelayID) > 0 {
		*sigpathIDs = append(*sigpathIDs, fx.DelayID)
	} else if len(fx.DistortionID) > 0 {
		*sigpathIDs = append(*sigpathIDs, fx.DistortionID)
	} else if len(fx.ModID) > 0 {
		*sigpathIDs = append(*sigpathIDs, fx.ModID)
	} else {
		*sigpathIDs = append(*sigpathIDs, fx.DspID)
	}
}

func getLicenseTier(data []byte) int {
	var s string
	for i, b := range data {
		if b == 0 {
			s = string(data[:i])
			break
		}
	}

	value, found := licenseTable[s]
	if !found {
		return 0
	}

	return value
}

func initFx2Map(filePath string, fx2Licenses map[string]int, packLicenses map[string]int) (err error) {
	var fxConfig BiasFxJSON
	fxConfig, err = readFxConfig(filePath)
	if err != nil {
		return
	}

	for _, cate := range fxConfig.BiasFx2 {
		for _, fx := range cate.FxList {
			data := LicenseDecrypt(fx.SigpathDescriptor, fx.DspID)
			fx2Licenses[fx.DspID] = getLicenseTier(data)

			initPack(fx, packLicenses)
		}
	}

	return
}

func initAmp2Map(filePath string, fx2Licenses map[string]int, amp2Licenses map[string]int, packLicenses map[string]int) (err error) {
	var fxConfig BiasFxJSON
	fxConfig, err = readFxConfig(filePath)
	if err != nil {
		return
	}

	for _, cate := range fxConfig.BiasFx2 {
		for _, fx := range cate.AmpList {
			data := LicenseDecrypt(fx.SigpathDescriptor, fx.DspID)
			fx2Licenses[fx.DspID] = getLicenseTier(data)

			initPack(fx, packLicenses)
		}
	}

	for _, cate := range fxConfig.BiasAMP2 {
		for _, fx := range cate.AmpList {
			data := LicenseDecrypt(fx.SigpathDescriptor, fx.DspID)
			amp2Licenses[fx.DspID] = getLicenseTier(data)

			initPack(fx, packLicenses)
		}
	}

	return
}

func initPack(fx FxStruct, packLicenses map[string]int) {
	if len(fx.PackID) > 0 || len(fx.PackID2) > 0 {
		merged := append([]string{}, append(fx.PackID, fx.PackID2...)...)
		value, found := extensionPack[merged[0]]
		if found {
			packLicenses[fx.DspID] = value
		}
	}
}

func readFxConfig(filePath string) (fxConfig BiasFxJSON, err error) {
	js, _ := OpenFile(filePath)
	err = json.Unmarshal(js, &fxConfig)
	return
}
