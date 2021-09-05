#!/usr/bin/env bash

#
# This is a utils shell script
#

# arguments: target, item1, item2, item3, ...
# returns 1 if target is in the given items, 0 otherwise.
util::array_contains() {
	local target="$1"
	local items
	shift
	for items; do
		if [[ "${items}" == "${target}" ]]; then
			return 1
		fi
	done
	return 0
}
