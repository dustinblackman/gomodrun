.DEFAULT_GOAL := build
VERSION := 0.0.1

build:
	cd cmd/gomodrun && go build -o ../../gomodrun -ldflags="-X main.version=$(VERSION)"

install:
	cd cmd/gomodrun && go install

lint:
	gomodrun golangci-lint run

lint-fix:
	gomodrun golangci-lint run --fix

release:
	gomodrun goreleaser

release-snapshot:
	gomodrun goreleaser --snapshot --skip-publish --rm-dist

test:
	gomodrun ginkgo -v -r .
	make build && rm gomodrun

test-coverage:
	gomodrun ginkgo -v -r -race -cover -coverprofile=coverage.txt -covermode=atomic -outputdir=. .
	make build && rm gomodrun
	gomodrun goveralls -coverprofile=coverage.txt -service=travis-ci -repotoken $$COVERALLS_TOKEN
