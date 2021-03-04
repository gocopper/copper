GO=GO111MODULE=on go
WIRE=wire

.PHONY: all
all: lint test

.PHONY: test
test:
	$(GO) test -cover ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: generate
generate:
	$(WIRE) .
