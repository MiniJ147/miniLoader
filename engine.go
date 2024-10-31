package main

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// combines runner and watcher together for a program life cycle
var runningCh = make(chan error)
var resetCh = make(chan bool, 1)

var activeTree = map[string]fs.FileInfo{}

type Engine struct {
	Cfg     *Config
	Watcher *Watcher
}

func CreateEngine() (*Engine, error) {
	cfg, err := SetupConfig()
	if err != nil {
		return nil, err
	}

	return &Engine{
		Cfg: cfg,
	}, nil
}

func (e *Engine) BuildProcess() error {
	fmt.Println("Building...")
	buildCmd := exec.Command("go", "build", "-o", e.Cfg.BuildDir, e.Cfg.MainFileArg)

	b, err := buildCmd.CombinedOutput()
	if err != nil {
		fmt.Println("ERROR: ", err, string(b))
		return err
	}

	fmt.Println("Build Result:", "success", string(b))
	return nil
}

func (e *Engine) RunProcess() context.CancelFunc {
	fmt.Println("Starting process...")
	ctx, cancel := context.WithCancel(context.Background())
	processCmd := exec.CommandContext(ctx, e.Cfg.BuildDir+"/main", e.Cfg.Args...)

	processCmd.Stdin = os.Stdin
	processCmd.Stdout = os.Stdout
	stderr, _ := processCmd.StderrPipe()

	go func() {
		err := processCmd.Start()
		if err != nil {
			runningCh <- err
		}

		b, _ := io.ReadAll(stderr)
		fmt.Println(string(b))

		runningCh <- processCmd.Wait()

	}()
	return cancel
}

func (e *Engine) Run() {
	err := e.BuildProcess()
	if err != nil {
		os.Exit(1)
	}

	go e.scanFiles(time.Second)
	cancelProcess := e.RunProcess()
	// scanner := bufio.NewScanner(stderr)

	for {
		select {
		case err := <-runningCh:
			if err == nil {
				return
			}

			fmt.Println("encountered error awating file change before continuing", err)

		case <-resetCh:
			fmt.Println("detected our reset")
			cancelProcess()

			err := e.BuildProcess()
			if err != nil {
				fmt.Println("failed to rebuild, waiting for file change to rebuild")
				continue
			}

			cancelProcess = e.RunProcess()
			runningCh = make(chan error)
		}
	}
}

func (e *Engine) scanFiles(timeoutDuration time.Duration) {
	for {
		shouldReset := false
		// fmt.Println(activeTree)

		filepath.WalkDir(e.Cfg.ProjectDir, func(path string, d fs.DirEntry, err error) error {
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
				// fmt.Println("file change detected now restarting program")
				shouldReset = true
				activeTree[path] = info
			}

			return nil
		})

		if shouldReset {
			resetCh <- true
		}

		time.Sleep(timeoutDuration)
	}
}

func (e *Engine) CleanUp() {
	fmt.Println("cleanning up...")
	e.Cfg.CleanUp()
}
