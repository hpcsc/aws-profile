#!/bin/bash

set -euo pipefail

go tool cover -html=coverage.txt -o coverage.html
