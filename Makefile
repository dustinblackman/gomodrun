.DEFAULT_GOAL := build
VERSION := 0.2.3

build:
	cd cmd/gomodrun && go build -o ../../gomodrun -ldflags="-X main.version=$(VERSION)"

install:
	cd cmd/gomodrun && go install -ldflags="-X main.version=$(VERSION)"

lint:
	gomodrun golangci-lint run

lint-fix:
	gomodrun golangci-lint run --fix

ginkgo:
	gomodrun ginkgo -v -r .

release:
	go mod tidy
	git add .
	git commit -m "v$(VERSION)"
	git tag -a "v$(VERSION)" -m "v$(VERSION)"
	git push
	git push --tags
	gomodrun goreleaser --rm-dist

release-snapshot:
	gomodrun goreleaser --snapshot --skip-publish --rm-dist

test: ginkgo
	@make build && rm gomodrun

test-coverage:
	gomodrun ginkgo -v -r -race -cover -coverprofile=coverage.txt -covermode=atomic -outputdir=. .
	@go tool cover -func=coverage.txt
	@make build && rm -f gomodrun
	@if [ "$$COVERALLS_TOKEN" != "" ]; then gomodrun goveralls -coverprofile=coverage.txt -service=travis-ci -repotoken $$COVERALLS_TOKEN; fi;
