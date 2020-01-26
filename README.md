# gomodrun

The forgotten go tool that executes and caches binaries included in go.mod files. This makes it easy to version cli tools in your projects such as `golangci-lint` and `ginkgo`.

## Example

```sh
  # Run a linter
  gomodrun golangci-lint run

  # Convert a JSON object to a Go struct, properly passing in stdin.
	echo example.json | gomodrun gojson > example.go
```

## Usage

gomodrun works by using a `tools.go` (or any other name) file that sits in the root of your project that contains all the CLI dependencies you want bundled in to your `go.mod`. Note the `// +build tools` at the top of the file is required, and allows you to name your tools file anything you like.

_tools.go_

```go
// +build tools

package myapp

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/onsi/ginkgo/ginkgo"
)
```

Run `go build tools.go` to add the dependencies to your `go.mod`. The build is expected to fail. Afterwards, you can run your tools by prefixing `gomodrun`. A binary will be built and cached in `.gomodrun` in the root of your project, allowing all runs after the first to be nice and fast.

```sh
  gomodrun golangci-lint run
```

## [License](./LICENSE)

MIT
