GO=GO111MODULE=on go
WIRE=wire

.PHONY: all
all: lint test

.PHONY: test
test:
	$(GO) test -tags csql_sqlite -cover ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: generate
generate:
	$(WIRE) .
