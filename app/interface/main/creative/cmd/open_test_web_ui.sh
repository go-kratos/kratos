#!/bin/bash

command -v goconvey >/dev/null 2>&1 || { echo >&2 "required goconvey but it's not installed."; echo "Aborting."; echo "Please run commond: go get github.com/smartystreets/goconvey"; exit 1; }

cd ../creative
goconvey -excludedDirs "vendor,node_modules,rpc" -packages 1
