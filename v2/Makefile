GO=GO111MODULE=on go

.PHONY: all
all: lint test

.PHONY: test
test:
	$(GO) test -cover ./...

.PHONY: lint
lint:
	golangci-lint run --fast --enable-all --exclude-use-default=false --fix

