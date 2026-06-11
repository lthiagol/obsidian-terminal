.PHONY: build run test lint vet clean install

BINARY := obsidian-terminal
GO := go
GOFLAGS := -ldflags="-s -w"

build:
	$(GO) build $(GOFLAGS) -o $(BINARY) .

run:
	$(GO) run .

test:
	$(GO) test ./... -v -count=1

test-race:
	$(GO) test ./... -race -count=1

vet:
	$(GO) vet ./...

lint:
	@which golangci-lint > /dev/null || (echo "golangci-lint not found: brew install golangci-lint"; exit 1)
	golangci-lint run ./...

fmt:
	$(GO) fmt ./...

clean:
	rm -f $(BINARY)

install:
	$(GO) install .
