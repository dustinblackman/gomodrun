// Package gomodrun is the forgotten go tool that executes and caches binaries included in go.mod files.
// This makes it easy to version cli tools in your projects such as `golangci-lint`
// and `ginkgo` that are versioned locked to what you specify in `go.mod`.
// Binaries are cached by go version and package version.
package gomodrun

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func cleanEmptyDirectory(root string) error {
	dirs := []string{}
	err := filepath.Walk(root, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			dirs = append(dirs, filePath)
		}

		return nil
	})

	if err != nil {
		return err
	}

	for i := len(dirs) - 1; i >= 0; i-- {
		files, readErr := ioutil.ReadDir(dirs[i])
		if readErr != nil {
			return err
		}

		if len(files) != 0 {
			continue
		}

		err = os.Remove(dirs[i])
		if err != nil {
			return err
		}
	}

	return err
}

func getAllBins(root string) ([]string, error) {
	binPaths := []string{}
	err := filepath.Walk(root, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			binPaths = append(binPaths, filePath)
		}

		return nil
	})

	return binPaths, err
}

// Tidy cleans .gomodrun of any outdated binaries.
func Tidy(pkgRoot string) error {
	var err error
	if pkgRoot == "" {
		pkgRoot, err = GetPkgRoot()
		if err != nil {
			return err
		}
	}

	gmrRoot := path.Join(pkgRoot, ".gomodrun")
	if _, err = os.Stat(gmrRoot); os.IsNotExist(err) {
		return nil
	}

	goVersion, err := getGoVersion()
	if err != nil {
		return err
	}

	gmrRootFiles, err := ioutil.ReadDir(gmrRoot)
	if err != nil {
		return err
	}

	for _, file := range gmrRootFiles {
		if file.Name() != goVersion {
			err = os.RemoveAll(path.Join(gmrRoot, file.Name()))
			if err != nil {
				return err
			}
		}
	}

	binPaths, err := getAllBins(gmrRoot)
	if err != nil {
		return err
	}

	if len(binPaths) == 0 {
		return nil
	}

	pkg, err := getToolsPkg(pkgRoot)
	if err != nil {
		return err
	}

	mod, err := getGoMod(pkgRoot)
	if err != nil {
		return err
	}

	versionedImports := []string{}
	for _, modulePath := range pkg.Imports {
		for importPath, version := range mod.Require {
			if strings.HasPrefix(modulePath, importPath) {
				versionedImports = append(versionedImports, importPath+"@"+version)
				break
			}
		}
	}

	for _, binPath := range binPaths {
		validBin := false
		for _, versionedImport := range versionedImports {
			if strings.Contains(binPath, versionedImport) {
				validBin = true
				break
			}
		}

		if !validBin {
			err = os.Remove(binPath)
			if err != nil {
				return err
			}
		}
	}

	err = cleanEmptyDirectory(pkgRoot)
	if err != nil {
		return err
	}

	return nil
}
