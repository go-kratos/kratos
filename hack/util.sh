#!/usr/bin/env bash

# This is a common util functions shell script

# arguments: target, item1, item2, item3, ...
# returns 0 if target is in the given items, 1 otherwise.
function util::array_contains() {
	local target="$1"
	shift
	local items="$*"
	for item in ${items[*]}; do
		if [[ "${item}" == "${target}" ]]; then
			return 0
		fi
	done
	return 1
}

# find all go mod path
# returns an array contains mod path
function util::find_modules() {
	find . -not \( \
		\( \
		-path './output' \
		-o -path './.git' \
		-o -path '*/third_party/*' \
		-o -path '*/vendor/*' \
		\) -prune \
		\) -name 'go.mod' -print0 | xargs -0 -I {} dirname {}
}
