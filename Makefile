GO=GO111MODULE=on go
WIRE=wire

.PHONY: all
all: lint generate test

.PHONY: test
test:
	$(GO) test -tags csql_sqlite -cover ./...

.PHONY: lint
lint: tidy
	golangci-lint run

.PHONY: tidy
tidy:
	$(GO) mod tidy

.PHONY: generate
generate:
	$(WIRE) .
