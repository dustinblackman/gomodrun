// +build tools

package gomodrun

import (
	_ "github.com/dustinblackman/go-hello-world-test-no-gomod/hello-world-no-gomod"
	_ "github.com/dustinblackman/go-hello-world-test/hello-world"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/goreleaser/goreleaser"
	_ "github.com/mattn/goveralls"
	_ "github.com/onsi/ginkgo/ginkgo"
)
