package main

// InputBankListJSON : input file
type InputBankListJSON struct {
	Author string `json:"author,omitempty"`
	Src    string `json:"src"`
	Dst    string `json:"dst"`
	Report string `json:"report"`
	Temp   string `json:"temp"`
	Banks  []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"banks"`
}
