#!/bin/bash

set -e

PATH_TO_EXECUTABLE=$1
if [ -z "${PATH_TO_EXECUTABLE}" ]; then
    PATH_TO_EXECUTABLE=bin/aws-profile
    echo "=== path to executable is not found, rebuilding to ${PATH_TO_EXECUTABLE}"
    go build -o bin/aws-profile
fi;

export TERM=xterm-256color

function setup() {
    echo "=============== Setup =================="
    rm -rf ./e2e/tmp && mkdir ./e2e/tmp
    cp ./e2e/config/* ./e2e/tmp
}

function test_set_from_credentials() {
    echo "========= Set from credentials ========="
    ./e2e/expect/set-from-credentials.exp ${PATH_TO_EXECUTABLE}

    local current_profile=$(AWS_SHARED_CREDENTIALS_FILE=./e2e/tmp/credentials AWS_CONFIG_FILE=./e2e/tmp/config ${PATH_TO_EXECUTABLE} get)
    local expected=credentials_profile_2
    if [ "${current_profile}" != "${expected}" ]; then
        echo "FAILED [set-from-credentials] Expected [${expected}], but got [${current_profile}]"
        exit 1
    fi;
}

function test_set_from_config() {
    echo "========== Set from config ============="
    ./e2e/expect/set-from-config.exp ${PATH_TO_EXECUTABLE}

    local current_profile=$(AWS_SHARED_CREDENTIALS_FILE=./e2e/tmp/credentials AWS_CONFIG_FILE=./e2e/tmp/config ${PATH_TO_EXECUTABLE} get)
    local expected='profile config_profile_1'
    if [ "${current_profile}" != "${expected}" ]; then
        echo "FAILED: [set-from-config] expected [${expected}], but got [${current_profile}]"
        exit 1
    fi;
}

function test_set_region() {
    echo "========== Set region ==================="
    ./e2e/expect/set-region.exp ${PATH_TO_EXECUTABLE}

    local current_region=$(AWS_SHARED_CREDENTIALS_FILE=./e2e/tmp/credentials AWS_CONFIG_FILE=./e2e/tmp/config ${PATH_TO_EXECUTABLE} get-region)
    local expected='us-east-1'
    if [ "${current_region}" != "${expected}" ]; then
        echo "FAILED: [set-region] expected [${expected}], but got [${current_region}]"
        exit 1
    fi;
}

function test_upgrade_to_stable() {
    echo "========= Upgrade to stable ============"
    # copy executable to separate file so that upgrade will not destroy original file
    local executable_to_be_upgrade=${PATH_TO_EXECUTABLE}-stable
    cp -vf ${PATH_TO_EXECUTABLE} ${executable_to_be_upgrade}

    ${executable_to_be_upgrade} upgrade
    local actual=$(${executable_to_be_upgrade} version)
    rm -f ${executable_to_be_upgrade}

    local not_expected='aws-profile (000)'
    if [ "${actual}" = "${not_expected}" ]; then
        echo "FAILED: [upgrade-to-stable] Expected version to be different from [${not_expected}], but got [${actual}]"
        exit 1
    fi
}

function test_upgrade_to_prerelease() {
    echo "======== Upgrade to prerelease ========="
    # copy executable to separate file so that upgrade will not destroy original file
    local executable_to_be_upgrade=${PATH_TO_EXECUTABLE}-prerelease
    cp -vf ${PATH_TO_EXECUTABLE} ${executable_to_be_upgrade}

    ${executable_to_be_upgrade} upgrade --prerelease
    local actual=$(${executable_to_be_upgrade} version)
    rm -f ${executable_to_be_upgrade}

    local not_expected='aws-profile (000)'
    if [ "${actual}" = "${not_expected}" ]; then
        echo "FAILED: [upgrade-to-stable] Expected version to be different from [${not_expected}], but got [${actual}]"
        exit 1
    fi
}

function teardown() {
    rm -rf ./e2e/tmp
}

if [ ! -z "${GITHUB_TOKEN}" ]; then
    echo "=== Github token is set"
fi

setup
test_set_from_credentials
test_set_from_config
test_set_region
test_upgrade_to_stable
test_upgrade_to_prerelease
teardown

echo "=== OK"
