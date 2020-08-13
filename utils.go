// Package gomodrun is the forgotten go tool that executes and caches binaries included in go.mod files.
// This makes it easy to version cli tools in your projects such as `golangci-lint`
// and `ginkgo` that are versioned locked to what you specify in `go.mod`.
// Binaries are cached by go version and package version.
package gomodrun

import (
	"go/build"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"

	"github.com/sirkon/goproxy/gomod"
)

func getGoMod(root string) (*gomod.Module, error) {
	gomodPath := path.Join(root, "go.mod")
	data, err := ioutil.ReadFile(gomodPath)
	if err != nil {
		return nil, err
	}
	mod, err := gomod.Parse(gomodPath, data)
	if err != nil {
		return nil, err
	}

	return mod, nil
}

func getToolsPkg(root string) (*build.Package, error) {
	importContext := build.Default
	importContext.BuildTags = []string{"tools"}
	return importContext.ImportDir(root, 0)
}

func getGoVersion() (string, error) {
	goVersionOutput, err := exec.Command("go", "version").Output()
	if err != nil {
		return "", err
	}

	return strings.Split(string(goVersionOutput), " ")[2], nil
}
