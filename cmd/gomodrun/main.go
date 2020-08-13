// Package main compiles the binary that uses the gomodrun library.
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"

	"github.com/dustinblackman/gomodrun"
)

var (
	version = "v0.0.0"
	commit  = "master"
	date    = "unknown"
)

func exitWithError(err error) {
	color.Red("gomodrun: " + err.Error())
	os.Exit(1)
}

func main() {
	if len(os.Args) <= 1 {
		exitWithError(errors.New("no binary name provided"))
	}

	if os.Args[1] == "--help" || os.Args[1] == "--version" {
		fmt.Printf(`gomodrun %s
Build Date: %s
https://github.com/dustinblackman/gomodrun/commit/%s

The forgotten go tool that executes and caches binaries included in go.mod files.

Usage:
	gomodrun [flags] cli-name [parameters]

Example:
	gomodrun golangci-lint run
	echo example.json | gomodrun gojson > example.go
	gomodrun -r ./alternative-tools-dir golangci-lint run

Flags:
  -r, --pkg-root string  Specify alternative root directory containing a go.mod and tools file. Defaults to walking up the file tree to locate go.mod.
  -t, --tidy  Cleans .gomodrun of any outdated binaries.`, version, date, commit)
		os.Exit(0)
	}

	cmdPosition := 1
	argsPosition := 2
	pkgRoot := ""

	skipNext := false
	for idx, entry := range os.Args {
		if skipNext {
			skipNext = false
			continue
		}

		if entry == "-t" || entry == "--tidy" {
			err := gomodrun.Tidy(pkgRoot)
			if err != nil {
				exitWithError(err)
			}
			os.Exit(0)
		}

		if entry == "-r" || entry == "--pkg-root" {
			pkgRoot = os.Args[idx+1]
			skipNext = true
			continue
		}

		if !strings.HasPrefix(entry, "-") {
			cmdPosition = idx
			argsPosition = idx + 1
			break
		}
	}

	exitCode, err := gomodrun.Run(os.Args[cmdPosition], os.Args[argsPosition:], &gomodrun.Options{
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Env:     os.Environ(),
		PkgRoot: pkgRoot,
	})

	if err != nil {
		exitWithError(err)
	}
	os.Exit(exitCode)
}
