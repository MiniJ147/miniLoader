package main

import (
	"fmt"
	"os"
)

// var activeTree = map[string]fs.FileInfo{}

// func GetRootDir(execCmd string) (string, string, error) {
// 	wd, err := os.Getwd()
// 	if err != nil {
// 		return "", "", err
// 	}

// 	currDir := filepath.Dir(filepath.Join(wd, execCmd))
// 	orgDir := currDir
// 	for {
// 		// fmt.Println("walking on", currDir)

// 		_, err = os.ReadDir(currDir)
// 		if err != nil {
// 			return "", "", err
// 		}

// 		_, err = os.Stat(currDir + "/go.mod")
// 		if err == nil {
// 			return currDir, orgDir, nil
// 		}
// 		currDir = currDir[:strings.LastIndex(currDir, "\\")]
// 	}
// }

// func scanFiles(root string) bool {
// 	shouldReset := false
// 	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
// 		if err != nil {
// 			return err
// 		}

// 		info, err := d.Info()
// 		if err != nil {
// 			return err
// 		}

// 		if !info.IsDir() && !strings.Contains(path, ".go") {
// 			return nil
// 		}

// 		intialState, ok := activeTree[path]
// 		if !ok || !intialState.ModTime().Equal(info.ModTime()) || info.Size() != intialState.Size() {
// 			fmt.Println("need to trigger reset")
// 			shouldReset = true
// 			activeTree[path] = info
// 		}

// 		return nil
// 	})

// 	return shouldReset
// }

func main() {
	cfg, err := SetupConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer cfg.CleanUp() // cleans program gracefully

	fmt.Println(cfg)

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	// exe := exec.CommandContext(ctx, "test/cmd/test")
	// buff, err := exe.Output()
	// fmt.Println(string(buff), err)
	// dirName, err := os.MkdirTemp("", "mini-loader-build")
	// if err != nil {
	// 	panic(err)
	// }
	// defer os.RemoveAll(dirName)

	// fmt.Println(os.Getwd())
	// fmt.Println(dirName)

	// if len(os.Args) < 2 {
	// 	fmt.Println("args not long enough, missing input file")
	// 	os.Exit(1)
	// }

	// //setting cfg data
	// rootDir, orgDir, err := GetRootDir(os.Args[1])
	// if err != nil {
	// 	panic(err)
	// }

	// cfg := Config{
	// 	// ExecCmd: string{"build"},
	// 	RootDir: rootDir,
	// 	OrgDir:  orgDir,
	// }

	// fmt.Println(cfg)
	// cmd := exec.Command("go", "build", "-o", dirName, "main.go")
	// cmd.Dir = cfg.OrgDir
	// b, e := cmd.CombinedOutput()
	// fmt.Println(string(b), e, cmd.Dir)
	// fmt.Println("built")

	// runCmd := exec.Command(fmt.Sprintf("%v/main", dirName)) // exec args,)
	// b, e = runCmd.CombinedOutput()
	// fmt.Println(string(b), e)

	// time.Sleep(10 * time.Second)
	// scanFiles(cfg.RootDir)
	// // fmt.Println(activeTree)

	// fmt.Println("attempting to execute", cfg.ExecCmd, cfg.RootDir)

	// //executing go project listed by user
	// // exec.CommandContext()

	// ctx, cancel := context.WithCancel(context.Background())
	// defer cancel()

	// cmd := exec.CommandContext(ctx, "go", cfg.ExecCmd...)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	// err = cmd.Start()
	// if err != nil {
	// 	panic(err)
	// }
	// time.Sleep(5 * time.Second)
	// cmd.Process.Signal(os.Interrupt)
	// fmt.Println()
}
