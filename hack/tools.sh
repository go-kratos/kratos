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

echo "$KRAOTS_HOME"

source "${KRAOTS_HOME}/hack/tools.sh"

LINTER=${KRAOTS_HOME}/bin/golangci-lint
LINTER_CONFIG=${KRAOTS_HOME}/.golangci.yml
FAILURE_FILE=${KRAOTS_HOME}/hack/.lintcheck_failures

function find_modules() {
	find . -not \( \
		\( \
		-path './output' \
		-o -path './.git' \
		-o -path '*/third_party/*' \
		-o -path '*/vendor/*' \
		\) -prune \
		\) -name 'go.mod' -print0 | xargs -0 -I {} dirname {}
}

all_modules=$(find_modules)
failing_modules=()
while IFS='' read -r line; do failing_modules+=("$line"); done < <(cat "$FAILURE_FILE")

# functions
# lint all mod
function lint() {
	for mod in $all_modules; do
		local in_failing
		in_failing="$(util::array_contains "$mod" "${failing_modules[*]}")"
		echo "$mod" "$in_failing"
		if [[ "$in_failing" -ne "1" ]]; then
			pushd "$mod" >/dev/null &&
				echo "golangci lint $(sed -n 1p go.mod | cut -d ' ' -f2)" &&
				eval "${LINTER} run --timeout=5m --config=${LINTER_CONFIG}"
			popd >/dev/null || exit
		fi
	done
}

# test all mod
function test() {
	for dir in $all_modules; do
		if [[ ! " ${failing_modules[*]} " =~ ${mod} ]]; then
			pushd "$dir" >/dev/null &&
				echo "go test $(sed -n 1p go.mod | cut -d ' ' -f2)" &&
				go test ./...
			popd >/dev/null || exit
		fi
	done
}

# try to fix all mod with golangci-lint
function fix() {
	for dir in $all_modules; do
		if [[ ! " ${failing_modules[*]} " =~ ${mod} ]]; then
			pushd "$dir" >/dev/null &&
				echo "golangci fix $(sed -n 1p go.mod | cut -d ' ' -f2)" &&
				eval "${LINTER} run -v --fix --timeout=5m --config=${LINTER_CONFIG}"
			popd >/dev/null || exit
		fi
	done
}

function tidy() {
	for dir in $all_modules; do
		if [[ ! " ${failing_modules[*]} " =~ ${mod} ]]; then
			pushd "$dir" >/dev/null &&
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
