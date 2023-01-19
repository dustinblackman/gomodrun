![gomodrun](assets/banner.jpg)

[![Build Status](https://img.shields.io/github/workflow/status/dustinblackman/gomodrun/Test?branch=master)](https://github.com/dustinblackman/gomodrun/actions)
[![Release](https://img.shields.io/github/v/release/dustinblackman/gomodrun)](https://github.com/dustinblackman/gomodrun/releases)
[![Coverage Status](https://coveralls.io/repos/github/dustinblackman/gomodrun/badge.svg?branch=master)](https://coveralls.io/github/dustinblackman/gomodrun?branch=master)
[![Go Report Card](http://goreportcard.com/badge/dustinblackman/gomodrun)](http://goreportcard.com/report/dustinblackman/gomodrun)
[![Godocs](https://godoc.org/github.com/dustinblackman/gomodrun?status.svg)](https://pkg.go.dev/github.com/dustinblackman/gomodrun?tab=doc)

The forgotten go tool that executes and caches binaries included in go.mod files. This makes it easy to version cli tools in your projects such as `golangci-lint` and `ginkgo` that are versioned locked to what you specify in `go.mod`. Binaries are cached by go version and package version.

## Example

```sh
  # Run a linter
  gomodrun golangci-lint run

  # Convert a JSON object to a Go struct, properly passing in stdin.
  echo example.json | gomodrun gojson > example.go

  # Specifiy alternative root directory containing a go.mod and tools file.
  gomodrun -r ./alternative-tools-dir golangci-lint run

  # Clean your .gomodrun folder of unused binaries.
  gomodrun --tidy
```

## Install

### MacOS

```sh
brew install dustinblackman/tab/gomodrun
```

### Debian / Ubuntu

```sh
curl -s https://apt.dustinblackman.com/KEY.gpg | apt-key add -
curl -s https://apt.dustinblackman.com/dustinblackman.list > /etc/apt/sources.list.d/dustinblackman.list
sudo apt-get update
sudo apt-get install gomodrun
```

### Fedora / CentOS

```sh
dnf config-manager --add-repo https://yum.dustinblackman.com/config.repo
dnf install gomodrun
```

### Nix

```sh
nix-env -f '<nixpkgs>' -iA nur.repos.dustinblackman.gomodrun
```

### Windows

__Chocolatey__

<!-- choco-install start -->
```sh
choco install gomodrun --version=0.4.5
```
<!-- choco-install end -->

__Scoop__

```sh
scoop bucket add dustinblackman https://github.com/dustinblackman/scoop-bucket.git
scoop install gomodrun
```

### Manual

Download the pre-compiled binaries and packages from the [releases page](https://github.com/dustinblackman/gomodrun/releases) and
copy to the desired location.

### Source

```sh
go get -u github.com/dustinblackman/gomodrun/cmd/gomodrun
```

## Usage

gomodrun works by using a `tools.go` (or any other name) file that sits in the root of your project that contains all the CLI dependencies you want bundled in to your `go.mod`. Note the `// +build tools` at the top of the file is required, and allows you to name your tools file anything you like.

__tools.go__

```go
// +build tools

package myapp

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/onsi/ginkgo/ginkgo"
)
```

Run `go build tools.go` to add the dependencies to your `go.mod`. The build is expected to fail.

### CLI

You can run your tools by prefixing `gomodrun`. A binary will be built and cached in `.gomodrun` in the root of your project, allowing all runs after the first to be nice and fast.

```sh
  gomodrun golangci-lint run
```

### Programmatically

You can also use `gomodrun` as a library.

```go
package main

import (
	"os"

	"github.com/dustinblackman/gomodrun"
)

func main() {
	exitCode, err := gomodrun.Run("golangci-lint", []string{"run"}, &gomodrun.Options{
		Stdin:   os.Stdin,
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
		Env:     os.Environ(),
		PkgRoot: "",
	})
}
```


## [License](./LICENSE)

MIT
