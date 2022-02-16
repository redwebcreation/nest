.PHONY: help
help: ## print this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {gsub("\\\\n",sprintf("\n%22c",""), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: install
install: ## install dev tools
	go install gotest.tools/gotestsum@latest
	go install honnef.co/go/tools/cmd/staticcheck@master
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.44.0

.PHONY: lint
lint: ## lint the code
	@echo "Running staticcheck..."
	staticcheck -checks all,-ST1000 ./...
	@echo "Running golangci-lint..."
	golangci-lint run ./...
	@echo "Done!"

.PHONY: test
test: ## run tests
	gotestsum -f testname ./... -- -count=1

tests: test

.PHONY: fmt
fmt: ## format the code
	@go list -f {{.Dir}} ./... | xargs gofmt -w -s -d


.PHONY: test-coverage
test-coverage: ## run test coverage
	go test -count=1 -cover -covermode=atomic -race ./...