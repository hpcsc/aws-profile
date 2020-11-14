#!/bin/bash

set -euo pipefail

# G107 is about providing url to http.Get() as taint input.
# This is necessary in our case and there's no clear way to solve it beside making url const (which is not possible in this case)
gosec -exclude=G107 ./...
