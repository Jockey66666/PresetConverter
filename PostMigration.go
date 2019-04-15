package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/thedevsaddam/gojsonq"
)

// PostMigration : post migration
func PostMigration(inputBanks InputBankListJSON, bankTable map[string]string, presetSlice []PresetSliceStruct) {
	// step 1. read fx2 bank.json
	bankJs := readFX2BankJSON(inputBanks.Dst)

	// step 2. create bank.json
	createBankJSON(inputBanks.Temp, bankTable, &bankJs)

	// step 3. create preset.json
	createPresetJSON(inputBanks.Temp, bankJs)

	// step 4. create report
	createReport(inputBanks.Report, presetSlice, bankTable)
}

func readFX2BankJSON(root string) BankListJSON {
	data, err := OpenFile(root + "/bank.json")

	if err != nil {
		os.Exit(ErrorOpenBankJSONFailed)
	}

	var bankJs BankListJSON
	json.Unmarshal(data, &bankJs)

	return bankJs
}

func createBankJSON(root string, bankTable map[string]string, bankJs *BankListJSON) {

	// find display order
	displayOrder := getDisplayOrder(bankJs.LiveBanks)

	var bankList []string
	filepath.Walk(root, GetSubFolderList(&bankList))
	for _, bankFolder := range bankList {

		bankName := bankTable[bankFolder]

		if len(bankFolder) != 36 || len(bankName) == 0 {
			fmt.Println("Delete", bankFolder, "it's not a bank.")
			os.RemoveAll(root + "/" + bankFolder)
			continue
		}

		b := LiveBanksJSON{
			BankFolder:   bankFolder,
			BankName:     bankName,
			DisplayOrder: displayOrder,
			IsFactory:    false,
		}

		bankJs.LiveBanks = append(bankJs.LiveBanks, b)
		displayOrder++
	}

	// rename duplicate bank name
	renameBank(bankJs)

	js, _ := json.MarshalIndent(bankJs, "", "    ")
	SaveFile(root+"/bank.json", js)
}

func getDisplayOrder(banks []LiveBanksJSON) int {

	displayOrder := 0
	for _, bank := range banks {
		if bank.DisplayOrder > displayOrder {
			displayOrder = bank.DisplayOrder
		}
	}
	displayOrder++

	return displayOrder
}

func rename(oldName string) func() string {
	newName := oldName
	r, _ := regexp.Compile(`_\d+$`)
	return func() string {
		index := r.FindStringIndex(newName)
		if len(index) > 1 {

			// get last two element
			s := index[len(index)-2:]

			prefix := newName[:s[0]]
			postfix := newName[s[0]+1 : s[1]]
			serial, _ := strconv.Atoi(postfix)
			serial++
			newName = prefix + "_" + strconv.Itoa(serial)
		} else {
			newName += "_1"
		}
		return newName
	}
}

func getNewName(oldName string, existNames map[string]int) string {

	newName := rename(oldName)
	var currentName string
	for {
		currentName = newName()
		_, found := existNames[currentName]
		if !found {
			break
		}
	}

	return currentName
}

// rename bank if necessary avoid crash issue
func renameBank(bankJs *BankListJSON) {

	m := make(map[string]int)
	for i, bank := range bankJs.LiveBanks {
		_, found := m[bank.BankName]
		if found {
			newName := getNewName(bank.BankName, m)
			fmt.Println("Duplicate bank name", bank.BankName, "rename it to ", newName)
			bankJs.LiveBanks[i].BankName = newName
		} else {
			m[bank.BankName] = i
		}
	}
}

func createPresetJSON(root string, bankJs BankListJSON) {

	for _, bank := range bankJs.LiveBanks {
		bankPath := root + "/" + bank.BankFolder
		if _, err := os.Stat(bankPath); os.IsNotExist(err) {
			// ignore fx2 bank
			continue
		}

		var presetJs PresetJSON

		var presetList []string
		filepath.Walk(bankPath, GetSubFolderList(&presetList))

		for i, presetFolder := range presetList {
			// get preset name
			meta := bankPath + "/" + presetFolder + "/" + "meta.json"
			res := gojsonq.New().File(meta).From("name").Get()

			p := LivePresetsJSON{
				DisplayOrder: i,
				IsFavorite:   false,
				PresetName:   res.(string),
				PresetUUID:   presetFolder,
			}

			presetJs.LivePresets = append(presetJs.LivePresets, p)
		}

		// save as preset.json
		js, _ := json.MarshalIndent(presetJs, "", "    ")
		SaveFile(bankPath+"/preset.json", js)
	}
}

func createReport(root string, presetSlice []PresetSliceStruct, bankTable map[string]string) {
	var reportJs ReportJSON

	// get file name
	t := time.Now()
	timeString := fmt.Sprintf("%d%02d%02d", t.Year(), t.Month(), t.Day())
	fileName := "MigrationReport_" + timeString + "_0"

	// rename by file list
	m := make(map[string]int)
	files, _ := ioutil.ReadDir(root) //specify the current dir
	for i, file := range files {
		if !file.IsDir() {
			fn := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			m[fn] = i
		}
	}
	fileName = getNewName(fileName, m)

	// set time
	reportJs.Time = t.Format(time.RFC3339Nano)

	// parse presetSlice
	getReportResult(presetSlice, &reportJs, bankTable)

	// save as MigrationReport_yyyymmdd_n.json
	js, _ := json.MarshalIndent(reportJs, "", "    ")
	SaveFile(root+"/"+fileName+".json", js)
}

func getReportResult(presetSlice []PresetSliceStruct, reportJs *ReportJSON, bankTable map[string]string) {

	m := make(map[string]map[string][]string)
	m["success"] = make(map[string][]string)
	m["failed"] = make(map[string][]string)

	// classification
	for _, preset := range presetSlice {
		if preset.MigrateResult > 0 {
			// success
			m["success"][preset.OldBankUUID] = append(m["success"][preset.OldBankUUID], preset.PresetName)
		} else {
			// failed
			m["failed"][preset.OldBankUUID] = append(m["failed"][preset.OldBankUUID], preset.PresetName)
		}
	}

	// add success presets
	reportJs.Success = make([]BanksReport, len(m["success"]))

	var count int
	count = 0
	for fx1BankUUID, presets := range m["success"] {
		var bankRP BanksReport
		bankRP.BankID = fx1BankUUID
		bankRP.BankName = bankTable[fx1BankUUID]

		bankRP.Presets = make([]string, len(presets))
		for i, p := range presets {
			bankRP.Presets[i] = p
		}

		reportJs.Success[count] = bankRP
		count++
	}

	// add failed presets
	reportJs.Failed = make([]BanksReport, len(m["failed"]))

	count = 0
	for fx1BankUUID, presets := range m["failed"] {
		var bankRP BanksReport
		bankRP.BankID = fx1BankUUID
		bankRP.BankName = bankTable[fx1BankUUID]

		bankRP.Presets = make([]string, len(presets))
		for i, p := range presets {
			bankRP.Presets[i] = p
		}

		reportJs.Failed[count] = bankRP
		count++
	}

}
