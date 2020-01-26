.DEFAULT_GOAL := build
VERSION := 0.0.1

build:
	cd cmd/gomodrun && go build -o ../../gomodrun -ldflags="-X main.version=$(VERSION)"

install:
	cd cmd/gomodrun && go install -ldflags="-X main.version=$(VERSION)"

lint:
	gomodrun golangci-lint run -v

lint-fix:
	gomodrun golangci-lint run --fix

# Add go build to test.
test:
	gomodrun ginkgo -v -r .
	make build && rm gomodrun
