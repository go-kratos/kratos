#!/bin/bash

set -ex

# Find all of our external imports and update them.
updated=false
for attempt in 1 2 3; do
  if ( go list -f '{{join .Imports "\n"}}' ./... && go list -f '{{join .TestImports "\n"}}' ./...; ) \
    | sort -u \
    | fgrep -v github.com/tsuna/gohbase \
    | xargs go get -d -f -u -v; then
    updated=true
    break
  fi
  sleep $((attempt*attempt))
done
if ! $updated; then
  echo failed to update dependencies
  exit 1
fi
