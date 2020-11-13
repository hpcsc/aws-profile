#!/bin/bash

set -euo pipefail

go test -coverprofile=coverage.txt -v ./...
