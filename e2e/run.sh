#!/bin/bash

set -e

PATH_TO_EXECUTABLE=$1
if [ -z "${PATH_TO_EXECUTABLE}" ]; then
    PATH_TO_EXECUTABLE=bin/aws-profile
    echo "=== path to executable is not found, rebuilding to ${PATH_TO_EXECUTABLE}"
    go build -o bin/aws-profile
fi;

export TERM=xterm-256color

rm -rf ./e2e/tmp && mkdir ./e2e/tmp
cp ./e2e/config/* ./e2e/tmp

./e2e/expect/set-from-credentials.exp ${PATH_TO_EXECUTABLE}

CURRENT_PROFILE=$(${PATH_TO_EXECUTABLE} get --credentials-path ./e2e/tmp/credentials --config-path ./e2e/tmp/config)
EXPECTED=credentials_profile_2
if [ "${CURRENT_PROFILE}" != "${EXPECTED}" ]; then
    echo "FAILED [set-from-credentials] Expected [${EXPECTED}], but got [${CURRENT_PROFILE}]"
    exit 1
fi;

./e2e/expect/set-from-config.exp ${PATH_TO_EXECUTABLE}

CURRENT_PROFILE=$(${PATH_TO_EXECUTABLE} get --credentials-path ./e2e/tmp/credentials --config-path ./e2e/tmp/config)
EXPECTED='assuming profile config_profile_1'
if [ "${CURRENT_PROFILE}" != "${EXPECTED}" ]; then
    echo "FAILED: [set-from-config] Expected [${EXPECTED}], but got [${CURRENT_PROFILE}]"
    exit 1
fi;

echo "=== OK"
rm -rf ./e2e/tmp
