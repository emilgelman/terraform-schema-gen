BIN_DIR := .tools/bin
GOLANGCI_LINT_VERSION := 1.31.0
GOLANGCI_LINT := $(BIN_DIR)/golangci-lint_$(GOLANGCI_LINT_VERSION)


all: run-test lint

run-test:
	go test -v ./...

lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run --fast

$(GOLANGCI_LINT):
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(BIN_DIR) v$(GOLANGCI_LINT_VERSION)
	mv $(BIN_DIR)/golangci-lint $(GOLANGCI_LINT)
