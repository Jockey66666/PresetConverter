package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/thedevsaddam/gojsonq"
	"github.com/tidwall/gjson"
)

// SplitAmpHeadCab : split amp to amp head and amp cab
func SplitAmpHeadCab(fx FxStruct, fx2Preset *Fx2PresetData) (err error) {

	if len(fx.Ampdata) > 0 {
		//embedded data
		var ampData EmbeddedAmpData
		ampData.ID = fx.Uniqueid
		ampData.AmpID = fx.AmpID
		ampData.EmbeddedType = "BiasAmp"

		var extraData interface{}
		err = json.Unmarshal([]byte(fx.Ampdata), &extraData)

		if err != nil {
			fmt.Println(err)
			return
		}

		m := extraData.(map[string]interface{})

		ampData.AmpData = m["ampData"]
		ampData.MetaData = m["metaData"]
		ampData.PanelData = m["panelData"]

		fx2Preset.Embedded = append(fx2Preset.Embedded, ampData)

		// cab sigpath
		var sigpathElement SigpathElement
		sigpathElement.ModulePresetName = ""
		sigpathElement.Active, err = strconv.ParseBool(fx.Active)

		if err != nil {
			return
		}

		sigpathElement.AmpType = "AmpCab"
		sigpathElement.AmpID = fx.AmpID
		sigpathElement.DspID = "BiasAmp"
		sigpathElement.ID = fx.Uniqueid
		sigpathElement.Selected, err = strconv.ParseBool(fx.Selected)

		if err != nil {
			return
		}

		ampDataByte, _ := json.Marshal(ampData.AmpData)
		query := gojsonq.New().JSONString(string(ampDataByte)).
			From("sigPath.blocks.items").Where("id", "=", "bias.cab").Only("params")
		var queryString []byte
		queryString, err = json.Marshal(query)

		if err != nil {
			return
		}

		result := gjson.Get(string(queryString), "0.params")
		result.ForEach(func(key, value gjson.Result) bool {
			r1 := gjson.Parse(value.String()).Get("id")
			r2 := gjson.Parse(value.String()).Get("value")

			var fx2Param SigpathParam
			fx2Param.ID = int(r1.Int()) // int64 to int
			fx2Param.Value = r2.Float()
			sigpathElement.Param = append(sigpathElement.Param, fx2Param)

			return true
		})

		fx2Preset.SigPath = append(fx2Preset.SigPath, sigpathElement)
	}

	return
}
