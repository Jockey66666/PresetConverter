package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/thedevsaddam/gojsonq"
	"github.com/tidwall/gjson"
)

// HandleAmp : Amp converter
func HandleAmp(fx FxElement, sigpath *[]SigpathElement, embedded *[]EmbeddedData) {

	sigpathElement := GetSigpathElement(fx)
	sigpathElement.AmpType = "AmpHead"
	sigpathElement.AmpID = fx.AmpID
	*sigpath = append(*sigpath, sigpathElement)

	err := splitAmpHeadCab(fx, sigpath, embedded)

	if err != nil {
		fmt.Println("convert amp error", err)
	}
}

func splitAmpHeadCab(fx FxElement, sigpath *[]SigpathElement, embedded *[]EmbeddedData) (err error) {

	if len(fx.Ampdata) > 0 {
		//embedded data
		var ampData EmbeddedData
		ampData.ID = fx.Uniqueid
		ampData.AmpID = fx.AmpID
		ampData.EmbeddedType = fx.Descriptor

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

		*embedded = append(*embedded, ampData)

		// cab sigpath
		var sigpathElement SigpathElement
		sigpathElement.Active, err = strconv.ParseBool(fx.Active)

		if err != nil {
			return
		}

		sigpathElement.AmpType = "AmpCab"
		sigpathElement.AmpID = fx.AmpID
		sigpathElement.DspID = fx.Descriptor
		sigpathElement.ID = fx.Uniqueid
		sigpathElement.Selected, err = strconv.ParseBool(fx.Selected)

		if err != nil {
			return
		}

		var ampDataByte []byte
		ampDataByte, err = json.Marshal(ampData.AmpData)

		if err != nil {
			return
		}

		ampDataString := string(ampDataByte)
		query := gojsonq.New().JSONString(ampDataString).
			From("sigPath.blocks.items").Where("id", "=", "bias.cab").
			OrWhere("id", "=", "bias.cab2").OrWhere("id", "=", "bias.cab2celestion").
			OrWhere("id", "=", "bias.IRLoader").Only("params")
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

		*sigpath = append(*sigpath, sigpathElement)
	}

	return
}
