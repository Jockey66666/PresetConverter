package main

import "strconv"

var treadplateCabParams = []float64{
	0.013000000268220901,
	0,
	0.75,
	0.33116883039474487,
	0.0010000000474974513,
	0.81746029853820801,
	0,
	0.38507461547851562,
	0.43382352590560913,
	0.10000000149011612,
	0,
	0.75,
	0,
	0.80000001192092896,
	0.5,
	0.10000000149011612,
}

var green25sCabParams = []float64{
	0.014000000432133675,
	0,
	0.75,
	0,
	0.0010000000474974513,
	0.75,
	0,
	0.43736007809638977,
	0.51470589637756348,
	0.05106707289814949,
	0,
	0.75,
	0,
	0.80000001192092896,
	0.5,
	0.10000000149011612,
}

var jazzCleanCabParams = []float64{
	0.0040000001899898052,
	0,
	0.75,
	0,
	0.0010000000474974513,
	0.75,
	0,
	0.39104476571083069,
	0.41911765933036804,
	0.10000000149011612,
	0,
	0.75,
	0,
	0.80000001192092896,
	0.5,
	0.10000000149011612,
}

// LoomisMetal : Loomis Metal
func LoomisMetal(fx FxElement, sigpath *[]SigpathElement) {
	handleMetalPack(fx, sigpath, treadplateCabParams)
}

// MerrowFire : Merrow Fire
func MerrowFire(fx FxElement, sigpath *[]SigpathElement) {
	handleMetalPack(fx, sigpath, green25sCabParams)
}

// Merrow5153 : Merrow 5153
func Merrow5153(fx FxElement, sigpath *[]SigpathElement) {
	handleMetalPack(fx, sigpath, green25sCabParams)
}

// OlaWar : Ola War
func OlaWar(fx FxElement, sigpath *[]SigpathElement) {
	handleMetalPack(fx, sigpath, treadplateCabParams)
}

// OlaPeace : Ola Peace
func OlaPeace(fx FxElement, sigpath *[]SigpathElement) {
	handleMetalPack(fx, sigpath, green25sCabParams)
}

// Loomis120 : Loomis 120
func Loomis120(fx FxElement, sigpath *[]SigpathElement) {
	handleMetalPack(fx, sigpath, jazzCleanCabParams)
}

func handleMetalPack(fx FxElement, sigpath *[]SigpathElement, cabParam []float64) {
	sigpathElement := GetSigpathElement(fx)
	sigpathElement.DspID = "FX2." + sigpathElement.DspID
	sigpathElement.AmpType = "PackAmpHead"

	if len(sigpathElement.Param) > 0 {
		start := len(sigpathElement.Param)
		end := len(sigpathElement.Param) + 4
		for i := start; i < end; i++ {
			p := SigpathParam{
				ID:    i,
				Value: 1.0,
			}
			sigpathElement.Param = append(sigpathElement.Param, p)
		}
	}

	*sigpath = append(*sigpath, sigpathElement)

	addCab(fx, sigpath, cabParam)
}

func addCab(fx FxElement, sigpath *[]SigpathElement, cabParam []float64) {
	var sigpathElement SigpathElement
	sigpathElement.AmpType = "PackAmpCab"
	sigpathElement.Active, _ = strconv.ParseBool(fx.Active)
	sigpathElement.DspID = "bias.cab2"
	sigpathElement.ID = fx.Uniqueid
	sigpathElement.Selected, _ = strconv.ParseBool(fx.Selected)

	sigpathElement.Param = make([]SigpathParam, len(cabParam))

	for i, param := range cabParam {
		p := SigpathParam{
			ID:    i,
			Value: param,
		}
		sigpathElement.Param[i] = p
	}

	*sigpath = append(*sigpath, sigpathElement)

}
