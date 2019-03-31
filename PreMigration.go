package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/tidwall/gjson"
)

// PreMigration : prepare to migrate
func PreMigration() (InputBankListJSON, []PresetSliceStruct, map[string]string) {

	// step 1. create report folder and delete temp folder
	CreateDirIfNotExist("report")
	CreateDirIfNotExist("temp")
	RemoveContents("temp")

	// step 2. read input_bank_list.json
	data, _ := OpenFile("input_bank_list.json")
	var inputBanks InputBankListJSON
	err := json.Unmarshal(data, &inputBanks)
	if err != nil {
		os.Exit(ErrorOpenBankListFailed)
	}

	var presetSlice []PresetSliceStruct
	var bankTable = make(map[string]string)
	// step 3. get preset list from banks
	for _, bank := range inputBanks.Banks {
		presetPath := inputBanks.Src + "/" + bank.ID + "/preset.json"
		data, _ := OpenFile(presetPath)
		var presetJs PresetJSON
		err := json.Unmarshal(data, &presetJs)
		if err != nil {
			continue
		}

		// create bank folder
		uuid := strings.ToUpper(uuid.Must(uuid.NewRandom()).String())
		CreateDirIfNotExist("temp/" + uuid)
		bankTable[uuid] = bank.Name

		for _, preset := range presetJs.LivePresets {
			p := PresetSliceStruct{
				BankUUID:   uuid,
				PresetName: preset.PresetName,
				PresetPath: inputBanks.Src + "/" + bank.ID + "/" + preset.PresetName + ".preset",
			}

			presetSlice = append(presetSlice, p)
		}
	}

	// step 4. clear dangling fx2 preset folder
	clearDanglingBank(inputBanks.Dst)

	return inputBanks, presetSlice, bankTable
}

func clearDanglingBank(root string) {

	// step 1. read bank.json and get LiveBanks
	data, err := OpenFile(root + "/bank.json")

	if err != nil {
		os.Exit(ErrorOpenBankJSONFailed)
	}

	js := string(data)

	// step 2. get bank folder list
	var bankList []string
	filepath.Walk(root, GetSubFolderList(&bankList))

	// step 3. remove the folder if it's not exist in bank.json
	for _, bankFolder := range bankList {
		result := gjson.Get(js, "LiveBanks.#[bank_folder="+bankFolder+"]")

		if !result.Exists() {
			fmt.Println("Delete", bankFolder)
			os.RemoveAll(root + "/" + bankFolder)
		}
	}
}
