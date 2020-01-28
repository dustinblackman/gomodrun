package main

import (
	"errors"
	"fmt"
	"os"

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
	gomodrun cli-name [parameters]

Example:
	gomodrun golangci-lint run
	echo example.json | gomodrun gojson > example.go`, version, date, commit)
		os.Exit(0)
	}

	exitCode, err := gomodrun.Run(os.Args[1], os.Args[2:], gomodrun.Options{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    os.Environ(),
	})
	if err != nil {
		exitWithError(err)
	}
	os.Exit(exitCode)
}
