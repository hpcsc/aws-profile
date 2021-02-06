#!/bin/bash

set -euo pipefail

jfrog rt upload \
        --access-token "${ACCESS_TOKEN}" \
        --url https://hpcsc.jfrog.io/artifactory \
        'batect*' \
        "${REPOSITORY}/${VERSION}/"
