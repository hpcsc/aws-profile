#!/bin/bash

set -euo pipefail

jfrog rt upload \
        --access-token "${ACCESS_TOKEN}" \
        --url https://hpcsc.jfrog.io/artifactory \
        "${SOURCE_PATTERN}" \
        "${REPOSITORY}/${VERSION_NUMBER}/"
