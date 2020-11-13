#!/bin/bash

set -euo pipefail

go vet ./...

golangci-lint run
