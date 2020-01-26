package gomodrun

import (
	"errors"
	"fmt"
	"go/build"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/sirkon/goproxy/gomod"
)

// Options contains parameters that are passed to `exec.Command` when running the binary.
type Options struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	Env    []string
}

// GetPkgRoot gets your projects package root, allowing you to run gomodrun from any sub directory.
func GetPkgRoot() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(path.Join(currentDir, "go.mod")); !os.IsNotExist(err) {
			absPath, err := filepath.Abs(currentDir)
			if absPath == "/" {
				return "", errors.New("go.mod not found")
			}
			if err != nil {
				return "", err
			}
			return absPath, nil
		}
		currentDir = path.Dir(currentDir)
	}
}

// GetCommandVersionedPkgPath extracts the command line tools package path and version from go.mod.
func GetCommandVersionedPkgPath(pkgRoot, binName string) (string, error) {
	importContext := build.Default
	importContext.BuildTags = []string{"tools"}
	pkg, err := importContext.ImportDir(pkgRoot, 0)
	if err != nil {
		return "", err
	}

	binModulePath := ""
	for _, modulePath := range pkg.Imports {
		if strings.HasSuffix(modulePath, binName) {
			binModulePath = modulePath
			break
		}
	}

	if binModulePath == "" {
		return "", errors.New("can not find bin path")
	}

	gomodPath := path.Join(pkgRoot, "go.mod")
	data, err := ioutil.ReadFile(gomodPath)
	if err != nil {
		return "", err
	}
	mod, err := gomod.Parse(gomodPath, data)
	if err != nil {
		return "", err
	}

	cmdPath := ""
	for importPath, version := range mod.Require {
		if strings.HasPrefix(binModulePath, importPath) {
			cmdPath = path.Join(importPath+"@"+version, strings.ReplaceAll(binModulePath, importPath, ""))
		}
	}

	if cmdPath == "" {
		return "", fmt.Errorf("can not find require for module %s in go.mod", binModulePath)
	}

	return cmdPath, nil
}

// GetCachedBin returns the path to the cached binary, building it if it doesn't exist.
func GetCachedBin(pkgRoot, binName, cmdPath string) (string, error) {
	goVersionOutput, err := exec.Command("go", "version").Output()
	if err != nil {
		return "", err
	}

	cachedBin, err := filepath.Abs(path.Join(pkgRoot, ".gomodrun/", strings.Split(string(goVersionOutput), " ")[2], cmdPath, binName))
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(cachedBin); os.IsNotExist(err) {
		moduleBinSrcPath := path.Join(os.Getenv("GOPATH"), "pkg/mod", cmdPath)
		if _, err := os.Stat(moduleBinSrcPath); os.IsNotExist(err) {
			return "", fmt.Errorf("module %s not downloaded. Run `go mod download`", cmdPath)
		}

		moduleSrcRoot := moduleBinSrcPath
		for {
			if strings.Contains(path.Base(moduleSrcRoot), "@") {
				break
			}
			moduleSrcRoot = path.Dir(moduleSrcRoot)
		}

		if _, err := os.Stat(path.Join(moduleSrcRoot, "go.mod")); os.IsNotExist(err) {
			pkgName := strings.Split(strings.Split(moduleSrcRoot, "pkg/mod/")[1], "@")[0]
			cmd := exec.Command("go", "mod", "init", pkgName)
			cmd.Dir = moduleSrcRoot
			output, err := cmd.CombinedOutput()
			if err != nil {
				return "", fmt.Errorf("initializing modules %s go.mod failed: %s", pkgName, output)
			}
		}

		err := os.MkdirAll(path.Dir(cachedBin), os.ModePerm)
		if err != nil {
			return "", err
		}

		cmd := exec.Command("go", "build", "-o", cachedBin)
		cmd.Dir = moduleBinSrcPath
		output, err := cmd.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("building %s failed: %s", binName, output)
		}
	}

	return cachedBin, nil
}

// Run executes your binary.
func Run(binName string, args []string, options Options) (int, error) {
	pkgRoot, err := GetPkgRoot()
	if err != nil {
		return -1, err
	}

	cmdPath, err := GetCommandVersionedPkgPath(pkgRoot, binName)
	if err != nil {
		return -1, err
	}

	cachedBin, err := GetCachedBin(pkgRoot, binName, cmdPath)
	if err != nil {
		return -1, err
	}

	cmd := exec.Command(cachedBin, args...)
	cmd.Stdin = options.Stdin
	cmd.Stderr = options.Stderr
	cmd.Stdout = options.Stdout
	cmd.Env = options.Env
	err = cmd.Run()
	if exiterr, ok := err.(*exec.ExitError); ok {
		if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
			return status.ExitStatus(), nil
		}
	}

	return 0, nil
}
