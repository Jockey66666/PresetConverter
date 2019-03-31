package main

// ReportJSON : report
type ReportJSON struct {
	Success []BanksReport `json:"success"`
	Failed  []BanksReport `json:"failed"`
	Time    string        `json:"time"`
}

// BanksReport : bank list
type BanksReport struct {
	BankID   string   `json:"bank_id"`
	BankName string   `json:"bank_name"`
	Presets  []string `json:"presets"`
}
