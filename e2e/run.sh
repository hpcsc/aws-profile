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

CURRENT_PROFILE=$(AWS_SHARED_CREDENTIALS_FILE=./e2e/tmp/credentials AWS_CONFIG_FILE=./e2e/tmp/config ${PATH_TO_EXECUTABLE} get)
EXPECTED=credentials_profile_2
if [ "${CURRENT_PROFILE}" != "${EXPECTED}" ]; then
    echo "FAILED [set-from-credentials] Expected [${EXPECTED}], but got [${CURRENT_PROFILE}]"
    exit 1
fi;

./e2e/expect/set-from-config.exp ${PATH_TO_EXECUTABLE}

CURRENT_PROFILE=$(AWS_SHARED_CREDENTIALS_FILE=./e2e/tmp/credentials AWS_CONFIG_FILE=./e2e/tmp/config ${PATH_TO_EXECUTABLE} get)
EXPECTED='profile config_profile_1'
if [ "${CURRENT_PROFILE}" != "${EXPECTED}" ]; then
    echo "FAILED: [set-from-config] Expected [${EXPECTED}], but got [${CURRENT_PROFILE}]"
    exit 1
fi;

echo "=== OK"
rm -rf ./e2e/tmp
