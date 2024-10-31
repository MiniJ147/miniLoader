package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const BUILD_DIR = "mini-loader-builds"
const MAIN_PATH_IDX = 0

type Config struct {
	Args        []string
	MainFileArg string
	ProjectDir  string
	BuildDir    string
	CallingDir  string
}

func SetupConfig() (*Config, error) {
	fail := func(msg string, err error) (*Config, error) {
		fmt.Println("setup config fail:", msg, err)
		return nil, err
	}
	config := &Config{}

	// set build directory
	buildDir, err := os.MkdirTemp("", BUILD_DIR)
	if err != nil {
		return fail("failed temp file build", err)
	}
	config.BuildDir = buildDir

	// grabbing args
	if len(os.Args) < 2 || !strings.Contains(os.Args[1], "main.go") {
		return fail("not enough args expected main.go file", nil)
	}
	if len(os.Args) > 2 {
		config.Args = os.Args[2:]
	}

	config.MainFileArg = os.Args[1]

	err = setProjectDirectories(config, config.MainFileArg)
	return config, err
}

func setProjectDirectories(cfg *Config, filePath string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	//setting working directory
	cfg.CallingDir = wd

	currDir := filepath.Dir(filepath.Join(wd, filePath))

	for {
		_, err = os.Stat(currDir + "/go.mod")

		// if we found our root with the go.mod
		if err == nil {
			cfg.ProjectDir = currDir
			return nil
		}

		// go up one more directory
		currDir = currDir[:strings.LastIndex(currDir, "\\")]
	}
}

func (cfg *Config) CleanUp() {
	os.RemoveAll(cfg.BuildDir)
}
