#!/usr/bin/env bash

# tools shell script

# set -o errexit
set -o nounset
set -o pipefail

GO111MODULE=on
KRAOTS_HOME=$(
	cd "$(dirname "${BASH_SOURCE[0]}")" &&
		cd .. &&
		pwd
)

all_package=$(find . -not -path "*/vendor/*" -not -path "*/examples/*" -type f -name go.mod -print0 | xargs -0 -I {} dirname {})

LINTER=${KRAOTS_HOME}/bin/golangci-lint
LINTER_CONFIG=${KRAOTS_HOME}/.golangci.yml

failure_file=${KRAOTS_HOME}/script/.check_failures

find_files() {
	find . -not \( \
		\( \
		-wholename './output' \
		-o -wholename './.git' \
		-o -wholename '*/third_party/*' \
		-o -wholename '*/vendor/*' \
		\) -prune \
		\) -name 'go.mod'
}

failing_packages=()
while IFS='' read -r line; do failing_packages+=("$line"); done < <(cat "$failure_file")

# lint all mod
function lint() {
	for dir in $all_package; do
		pushd "$dir" >/dev/null &&
			echo "golangci lint $(sed -n 1p go.mod | cut -d ' ' -f2)" &&
			eval "${LINTER} run --timeout=5m --config=${LINTER_CONFIG}"
		popd >/dev/null
	done
}

# test all mod
function test() {
	for dir in $all_package; do
		pushd "$dir" >/dev/null &&
			echo "go test $(sed -n 1p go.mod | cut -d ' ' -f2)" &&
			go test ./...
		popd >/dev/null
	done
}

# try to fix all mod with golangci-lint
function fix() {
	for dir in $all_package; do
		pushd "$dir" >/dev/null &&
			echo "golangci fix $(sed -n 1p go.mod | cut -d ' ' -f2)" &&
			eval "${LINTER} run -v --fix --timeout=5m --config=${LINTER_CONFIG}"
		popd >/dev/null
	done
}

function tidy() {
	for dir in $all_package; do
		pushd "$dir" >/dev/null &&
			echo "go mod tidy $(sed -n 1p go.mod | cut -d ' ' -f2)" &&
			go mod tidy
		popd >/dev/null
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
