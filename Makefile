export GO111MODULE=on

PROJECT := $(CURDIR)
APP = randomweighttable
PKG = ...
LINT = $(shell golangci-lint --version | grep -Eo 'version [0-9]\.[0-9]*\.[0-9]*')

.PHONY: clean
clean:
	@rm -rf $(PROJECT)/bin coverage coverage.html
	@rm -rf $(PROJECT)/bin linter linter.yml

.PHONY: init
init: clean
	@mkdir -p $(PROJECT)/bin
	@go mod tidy

.PHONY: test
test: lint
	@mkdir -p $(PROJECT)/bin/coverage
	@go test -p 10 -coverprofile=bin/coverage/coverage.out ./$(PKG)
	@go tool cover -func=bin/coverage/coverage.out
	@go tool cover -html=bin/coverage/coverage.out -o bin/coverage/coverage.html

.PHONY: lint
lint: clean
	@echo ">> Linting Go code... $(LINT)"
	@golangci-lint run

.PHONY: build
build: init lint
	@echo '::build::'
	@go build -o $(PROJECT)/bin/$(APP) $(PROJECT)/cmd/$(APP)