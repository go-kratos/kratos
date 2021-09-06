#!/usr/bin/env bash

#
# This is a tools shell script
#

# set -o errexit
set -o nounset
set -o pipefail

GO111MODULE=on
KRAOTS_HOME=$(
	cd "$(dirname "${BASH_SOURCE[0]}")" &&
		cd .. &&
		pwd
)

source "${KRAOTS_HOME}/hack/util.sh"

LINTER=${KRAOTS_HOME}/bin/golangci-lint
LINTER_CONFIG=${KRAOTS_HOME}/.golangci.yml
FAILURE_FILE=${KRAOTS_HOME}/hack/.lintcheck_failures

all_modules=$(util::find_modules)
failing_modules=()
while IFS='' read -r line; do failing_modules+=("$line"); done < <(cat "$FAILURE_FILE")

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
		util::array_contains "$mod" "${failing_modules[*]}" && in_failing=$? || in_failing=$?
		if [[ "$in_failing" -ne "0" ]]; then
			pushd "$mod" >/dev/null &&
				echo "go test $(sed -n 1p go.mod | cut -d ' ' -f2)" &&
				go test ./...
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
		local in_failing
		util::array_contains "$mod" "${failing_modules[*]}" && in_failing=$? || in_failing=$?
		if [[ "$in_failing" -ne "0" ]]; then
			pushd "$mod" >/dev/null &&
				echo "go mod tidy $(sed -n 1p go.mod | cut -d ' ' -f2)" &&
				go mod tidy
			popd >/dev/null || exit
		fi
	done
}

function help() {
	echo "use: lint, test, fix, tidy"
}

case $1 in
lint)
	shift
	lint "$@"
	;;
test)
	shift
	test "$@"
	;;
tidy)
	shift
	tidy "$@"
	;;
fix)
	shift
	fix "$@"
	;;
*)
	help
	;;
esac
