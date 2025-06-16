#!/bin/bash
set -eu -o pipefail

GOMOCK_VERSION=0.5.1

go install "go.uber.org/mock/mockgen@v${GOMOCK_VERSION}"
