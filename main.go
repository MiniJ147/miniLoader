package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

/*
miniLoader - live reloader written in go
1. take input to the main.go file
2. return up until you find the go.mod file (this will be your root directory being watched)
3. grab instance of all .go files
4. start the input program (main.go)
5. start the watcher

/*
-W means getting watched

Approach 1: watch all files form main.go directory
	/cmd
		/app
			/extra - W
			main.go - W
	/internal
		... (go files)

Approach 2: watch files from executing directory
(wherever the terminal is calling from)
/*
	/cmd
		/app
			/extra - W
			main.go - W
	/internal
		... (go files)

terminal: $user/programs/example1> go run cmd/app/main.go
1. take user/programs/example1 - use this as the watching directory

ISSUE:
terminal: $user/programs> go run example1/cmd/app/main.go
it will watach all sub folders in programs/
EX: it will watch outside files not included in our project


Approach 3: .config file

Approach 4: find the .mod file

// Watch for
- Changes on existing file
- If new files are added
- if files are removed
*/

// nodemon server.js == node server.js
// go run main.go [watcher] /main
// watcher ./main.go

type Config struct {
	ExecCmd []string
	RootDir string
}

var activeTree = map[string]fs.FileInfo{}

func GetRootDir(execCmd string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	currDir := filepath.Dir(filepath.Join(wd, execCmd))
	for {
		// fmt.Println("walking on", currDir)

		_, err = os.ReadDir(currDir)
		if err != nil {
			return "", err
		}

		_, err = os.Stat(currDir + "/go.mod")
		if err == nil {
			return currDir, nil
		}
		currDir = currDir[:strings.LastIndex(currDir, "\\")]
	}
}

func scanFiles(root string) bool {
	shouldReset := false
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if !info.IsDir() && !strings.Contains(path, ".go") {
			return nil
		}

		intialState, ok := activeTree[path]
		if !ok || !intialState.ModTime().Equal(info.ModTime()) || info.Size() != intialState.Size() {
			fmt.Println("need to trigger reset")
			shouldReset = true
			activeTree[path] = info
		}

		return nil
	})

	return shouldReset
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("args not long enough, missing input file")
		os.Exit(1)
	}

	//setting cfg data
	rootDir, err := GetRootDir(os.Args[1])
	if err != nil {
		panic(err)
	}

	cfg := Config{
		ExecCmd: append([]string{"run"}, os.Args[1:]...),
		RootDir: rootDir,
	}

	scanFiles(cfg.RootDir)
	// fmt.Println(activeTree)

	fmt.Println("attempting to execute", cfg.ExecCmd, cfg.RootDir)

	//executing go project listed by user
	// exec.CommandContext()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", cfg.ExecCmd...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Start()
	if err != nil {
		panic(err)
	}
	go func() {
		time.Sleep(5 * time.Second)
		cmd.Cancel()
	}()
	cmd.Wait()
	fmt.Println("hello")
	// fmt.Println()
}
