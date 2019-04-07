package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/mitchellh/go-ps"
)

// ProcessGuard : only one instance can be executed.
type ProcessGuard struct {
	Args []string
}

// Guard : lock process by 'guard' arg
func (pg ProcessGuard) Guard() {

	// lock it forever
	if len(pg.Args) > 1 {
		last := pg.Args[len(pg.Args)-1]
		if last == "guard" {
			for {
				fmt.Println("guarding")
				time.Sleep(time.Minute)
			}
		}
	}

	pname := getCurrentProcessName(pg.Args[0])
	if pids := findProcessCount(pname); len(pids) > 0 {
		fmt.Println(pname, ", pid:", pids[0], ", already running")
		os.Exit(0)
	}
}

func getCurrentProcessName(path string) string {

	var slash string
	if runtime.GOOS == "windows" {
		slash = "\\"
	} else {
		slash = "/"
	}

	lastSlash := strings.LastIndex(path, slash)
	if lastSlash != -1 {
		return path[lastSlash+1 : len(path)]
	}

	return path
}

func findProcessCount(name string) []int {

	pid := []int{}
	ps, _ := ps.Processes()
	for i := range ps {
		if ps[i].Pid() != os.Getpid() && strings.EqualFold(ps[i].Executable(), name) {
			pid = append(pid, ps[i].Pid())
		}
	}

	return pid
}
