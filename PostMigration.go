package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// PostMigration : post migration
func PostMigration(fx2PresetPath string, bankTable map[string]string) {
	// step 1. read fx2 bank.json
	bankJs := readFX2BankJSON(fx2PresetPath)

	// step 2. create bank.json
	createBankJSON("temp", bankTable, bankJs)

	// step 3. cereate preset.json
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
			prefix := newName[:index[0]]
			postfix := newName[index[0]+1 : index[1]]
			serial, _ := strconv.Atoi(postfix)
			serial++
			newName = prefix + "_" + strconv.Itoa(serial)
		} else {
			newName += "_1"
		}
		return newName
	}
}

func getNewName(oldName string, existNames map[string]string) string {

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

	m := make(map[string]string)
	for i, bank := range bankJs.LiveBanks {
		_, found := m[bank.BankName]
		if found {
			newName := getNewName(bank.BankName, m)
			fmt.Println("Duplicate bank name", bank.BankName, "rename it to ", newName)
			bankJs.LiveBanks[i].BankName = newName
		} else {
			m[bank.BankName] = bank.BankName
		}
	}
}

func createBankJSON(root string, bankTable map[string]string, bankJs BankListJSON) {

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

	renameBank(&bankJs)
	js, _ := json.MarshalIndent(bankJs, "", "    ")
	SaveFile(root+"/bank.json", js)
}
