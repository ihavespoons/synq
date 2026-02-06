VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: build install clean test lint release-dry-run

build:
	go build $(LDFLAGS) -o synq ./cmd/synq

install:
	go install $(LDFLAGS) ./cmd/synq

clean:
	rm -f synq
	rm -rf dist/

test:
	go test -v ./...

lint:
	go vet ./...
	golangci-lint run ./...

release-dry-run:
	goreleaser release --snapshot --clean
