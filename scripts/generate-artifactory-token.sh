#!/bin/bash

set -euo pipefail

jfrog rt access-token-create "${NAME}" \
        --access-token "${ACCESS_TOKEN}" \
        --expiry 0 \
        --url https://hpcsc.jfrog.io/artifactory \
        --groups "${GROUP}"
