#!/bin/bash

# This is used by the linter action.
# Recursively finds all directories with a go.mod file and creates
# a GitHub Actions JSON output option.

echo "Resolving modules in $(pwd)"

PATHS=$(find . -not -path "*/vendor/*" -not -path "*/examples/*" -type f -name go.mod -printf '{"workdir":"%h"},')

echo "::set-output name=matrix::{\"include\":[${PATHS%?}]}"
