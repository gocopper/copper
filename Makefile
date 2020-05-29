GO=GO111MODULE=on go
GOIMPORTS=goimports

.PHONY: test
test:
	$(GO) test ./pkg/...

.PHONY: imports
imports:
	$(GOIMPORTS) -w .

