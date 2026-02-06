VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-s -w -X main.version=$(VERSION)"

.PHONY: build install clean test lint

build:
	go build $(LDFLAGS) -o synq ./cmd/synq

install:
	go install $(LDFLAGS) ./cmd/synq

clean:
	rm -f synq

test:
	go test ./...

lint:
	golangci-lint run
