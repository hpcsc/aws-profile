#!/bin/bash

set -euo pipefail

go build -i \
         -ldflags="-X github.com/hpcsc/aws-profile/internal/version.version=${VERSION}" \
         -o "${OUTPUT}" \
         ./cmd
