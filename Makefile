GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/dist
GOENVVARS := GOBIN=$(GOBIN)
GOBINARY := cosmos-utility
GOCMD := $(GOBASE)/cmd

LINT := $$(go env GOPATH)/bin/golangci-lint run --timeout=5m -E whitespace -E misspell -E gofmt -E goimports --exclude-use-default=false --max-same-issues 0
BUILD := $(GOENVVARS) go build $(LDFLAGS) -o $(GOBIN)/$(GOBINARY) $(GOCMD)

.PHONY: build
build: ## Builds the binary locally into ./dist
	$(BUILD)

.PHONY: lint
lint: ## Runs the linter
	$(LINT)

.PHONY: validator-status
validator-status: build ## Analyze the validator status
	./dist/cosmos-utility validator-status -c odin -p odin

.PHONY: vesting-analyze
vesting-analyze: build ## Analyze the validator status
	./dist/cosmos-utility vesting-analyze


.PHONY: test
test: ## Runs only short tests without checking race conditions
	go test --cover -short -p 1 ./...

## Help display.
## Pulls comments from beside commands and prints a nicely formatted
## display with the commands and their usage information.
.DEFAULT_GOAL := help

.PHONY: help
help: ## Prints this help
		@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'