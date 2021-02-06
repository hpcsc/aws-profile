#!/bin/bash

set -euo pipefail

jfrog rt download \
        --access-token "${ACCESS_TOKEN}" \
        --url https://hpcsc.jfrog.io/artifactory \
        --flat \
        "${REPOSITORY}/${VERSION_NUMBER}/*" \
        "${DESTINATION}"
