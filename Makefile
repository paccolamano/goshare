VERSION = $(shell scripts/git-version.sh)

.PHONY: all generate test lint fmt tidy version
all: test lint fmt

generate:
	scripts/install.sh mockgen
	find . -name go.mod -execdir go generate ./... \;

test: generate
	find . -name go.mod -execdir go test ./... -cover -count 1 \;

lint:
	scripts/install.sh golangci-lint
	find . -name go.mod -execdir golangci-lint run \;

fmt:
	find . -name go.mod -execdir go fmt ./... \;

tidy:
	find . -name go.mod -execdir go mod tidy \;

version:
	@echo $(VERSION)
