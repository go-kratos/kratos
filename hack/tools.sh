#!/usr/bin/env bash

# This is a tools shell script
# used by Makefile commands

set -o errexit
set -o nounset
set -o pipefail

GO111MODULE=on
KRATOS_HOME=$(
	cd "$(dirname "${BASH_SOURCE[0]}")" &&
		cd .. &&
		pwd
)

source "${KRATOS_HOME}/hack/util.sh"

LINTER=${KRATOS_HOME}/bin/golangci-lint
LINTER_CONFIG=${KRATOS_HOME}/.golangci.yml
FAILURE_FILE=${KRATOS_HOME}/hack/.lintcheck_failures
IGNORED_FILE=${KRATOS_HOME}/hack/.test_ignored_files

all_modules=$(util::find_modules)
failing_modules=()
while IFS='' read -r line; do failing_modules+=("$line"); done < <(cat "$FAILURE_FILE")
ignored_modules=()
while IFS='' read -r line; do ignored_modules+=("$line"); done < <(cat "$IGNORED_FILE")

# functions
# lint all mod
function lint() {
	for mod in $all_modules; do
		local in_failing
		util::array_contains "$mod" "${failing_modules[*]}" && in_failing=$? || in_failing=$?
		if [[ "$in_failing" -ne "0" ]]; then
			pushd "$mod" >/dev/null &&
				echo "golangci lint $(sed -n 1p go.mod | cut -d ' ' -f2)" &&
				eval "${LINTER} run --timeout=5m --config=${LINTER_CONFIG}"
			popd >/dev/null || exit
		fi
	done
}

# test all mod
function test() {
	for mod in $all_modules; do
		local in_failing
		util::array_contains "$mod" "${ignored_modules[*]}" && in_failing=$? || in_failing=$?
		if [[ "$in_failing" -ne "0" ]]; then
			pushd "$mod" >/dev/null &&
				echo "go test $(sed -n 1p go.mod | cut -d ' ' -f2)" &&
				go test -race ./...
			popd >/dev/null || exit
		fi
	done
}

function test_coverage() {
	echo "" >coverage.out
	local base
	base=$(pwd)
	for mod in $all_modules; do
		local in_failing
		util::array_contains "$mod" "${ignored_modules[*]}" && in_failing=$? || in_failing=$?
		if [[ "$in_failing" -ne "0" ]]; then
			pushd "$mod" >/dev/null &&
				echo "go test $(sed -n 1p go.mod | cut -d ' ' -f2)" &&
				go test -race -coverprofile=profile.out -covermode=atomic ./...
			if [ -f profile.out ]; then
				cat profile.out >>"${base}/coverage.out"
				rm profile.out
			fi
			popd >/dev/null || exit
		fi
	done
}

# try to fix all mod with golangci-lint
function fix() {
	for mod in $all_modules; do
		local in_failing
		util::array_contains "$mod" "${failing_modules[*]}" && in_failing=$? || in_failing=$?
		if [[ "$in_failing" -ne "0" ]]; then
			pushd "$mod" >/dev/null &&
				echo "golangci fix $(sed -n 1p go.mod | cut -d ' ' -f2)" &&
				eval "${LINTER} run -v --fix --timeout=5m --config=${LINTER_CONFIG}"
			popd >/dev/null || exit
		fi
	done
}

function tidy() {
	for mod in $all_modules; do
		pushd "$mod" >/dev/null &&
			echo "go mod tidy $(sed -n 1p go.mod | cut -d ' ' -f2)" &&
			go mod tidy
		popd >/dev/null || exit
	done
}

function help() {
	echo "use: lint, test, test_coverage, fix, tidy"
}

case $1 in
lint)
	lint
	;;
test)
	test
	;;
test_coverage)
	test_coverage
	;;
tidy)
	tidy
	;;
fix)
	fix
	;;
*)
	help
	;;
esac
