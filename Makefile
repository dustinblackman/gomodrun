.DEFAULT_GOAL := build

build:
	cd cmd/gomodrun && go build -o ../../gomodrun -ldflags="-X main.version=$$(cat ../../VERSION)"

install:
	cd cmd/gomodrun && go install -ldflags="-X main.version=$$(cat ../../VERSION)"

lint:
	gomodrun golangci-lint run

lint-fix:
	gomodrun golangci-lint run --fix

ginkgo:
	gomodrun ginkgo -v -r .

release:
	go mod tidy
	git add .
	git commit -m "v$$(cat VERSION)"
	git tag -a "v$$(cat VERSION)" -m "v$$(cat VERSION)"
	git push
	git push --tags
	gomodrun goreleaser --rm-dist
	.github/apt.sh "$$(cat VERSION)" "$$(realpath dist)"
	.github/nur.sh "$$(cat VERSION)" "$$(realpath dist)"

release-snapshot:
	gomodrun goreleaser --snapshot --skip-publish --rm-dist

test: ginkgo
	@make build && rm gomodrun

test-coverage:
	gomodrun ginkgo -v -r -cover -coverprofile=coverage.txt -covermode=atomic -outputdir=. .
	@go tool cover -func=coverage.txt
	@make build && rm -f gomodrun
