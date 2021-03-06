package main

/*

```

├── GlobalPresets
|   ├── bank.json
|   └── (banks)
└── MigrationTool
    ├── report
    |    └── MigrationReport_MMDD_N.json
    ├── temp
    |    ├── bank.json (output)
    |    └── (banks) (output)
    |
    ├── input_bank_list.json (input)
    ├── PresetConverterWin.exe
    └── PresetConverterMac

```

*/

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

func main() {
	t := time.Now()
	defer func() {
		elapsed := time.Since(t)
		fmt.Println("")
		fmt.Println("elapsed:", elapsed)
	}()

	// process guard, only one instance can be executed.
	pg := ProcessGuard{Args: os.Args}
	pg.Guard()

	// pre migration
	inputBanks, presetSlice, bankTable := PreMigration()
	var checker LicenseChecker
	checker.Init(inputBanks.Dst + "/../")

	// migration
	cpu := runtime.NumCPU()
	num := len(presetSlice)

	countChannel := make(chan int)
	for i := 0; i < cpu; i++ {
		from := i * num / cpu
		to := (i + 1) * num / cpu
		go func() {
			countChannel <- MigrationCore(checker, inputBanks, presetSlice[from:to])
		}()
	}

	total := 0
	for i := 0; i < cpu; i++ {
		total += <-countChannel
	}

	fmt.Println(total, "presets migrated")

	// post migration
	PostMigration(inputBanks, bankTable, presetSlice)
}
