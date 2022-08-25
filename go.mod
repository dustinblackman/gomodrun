module github.com/dustinblackman/gomodrun

go 1.18

require (
	github.com/dustinblackman/go-hello-world-test v0.0.2
	github.com/dustinblackman/go-hello-world-test-no-gomod v0.0.2
	github.com/fatih/color v1.10.0
	github.com/golangci/golangci-lint v1.35.2
	github.com/goreleaser/goreleaser v0.155.0
	github.com/mattn/goveralls v0.0.2
	github.com/novln/macchiato v1.0.1
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/otiai10/copy v1.0.2
	github.com/sirkon/goproxy v1.4.8
	mvdan.cc/gofumpt v0.1.0
)

replace github.com/novln/macchiato => github.com/dustinblackman/macchiato v0.0.0-20200814125024-987bc68e2aec
