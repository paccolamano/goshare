#!/bin/bash
set -eu -o pipefail

GOLANGCI_LINT_VERSION=2.0.1

go install "github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v${GOLANGCI_LINT_VERSION}"
