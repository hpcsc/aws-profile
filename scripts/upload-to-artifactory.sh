#!/bin/bash

set -euo pipefail

jfrog rt upload \
        --access-token "${ACCESS_TOKEN}" \
        --url https://hpcsc.jfrog.io/artifactory \
        "${SOURCE_PATTERN}" \
        "${REPOSITORY}/${VERSION_NUMBER}/"

jfrog rt sp \
        --recursive \
        --include-dirs \
        --access-token "${ACCESS_TOKEN}" \
        --url https://hpcsc.jfrog.io/artifactory \
        "${REPOSITORY}/${VERSION_NUMBER}*" \
        "version=${VERSION_NUMBER}"

rm -f /tmp/latest-version && echo "${VERSION_NUMBER}" > /tmp/latest-version
jfrog rt upload \
        --access-token "${ACCESS_TOKEN}" \
        --url https://hpcsc.jfrog.io/artifactory \
        /tmp/latest-version \
        "${REPOSITORY}/latest-version"
