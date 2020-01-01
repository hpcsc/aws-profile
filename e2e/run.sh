#!/bin/bash

set -e

export TERM=xterm-256color

go build

rm -rf ./e2e/tmp && mkdir ./e2e/tmp
cp ./e2e/config/* ./e2e/tmp

./e2e/set-from-credentials.exp

CURRENT_PROFILE=$(./aws-profile get --credentials-path ./e2e/tmp/credentials --config-path ./e2e/tmp/config)
EXPECTED=credentials_profile_2
if [ "${CURRENT_PROFILE}" != "${EXPECTED}" ]; then
    echo "FAILED [set-from-credentials] Expected [${EXPECTED}], but got [${CURRENT_PROFILE}]"
    exit 1
fi;

./e2e/set-from-config.exp

CURRENT_PROFILE=$(./aws-profile get --credentials-path ./e2e/tmp/credentials --config-path ./e2e/tmp/config)
EXPECTED='assuming profile config_profile_1'
if [ "${CURRENT_PROFILE}" != "${EXPECTED}" ]; then
    echo "FAILED: [set-from-config] Expected [${EXPECTED}], but got [${CURRENT_PROFILE}]"
    exit 1
fi;

echo "=== OK"
rm -rf ./e2e/tmp
