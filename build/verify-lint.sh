#!/bin/bash
# Copyright 2017 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o nounset
set -o pipefail

bgr -type=dir  -script=./app/tool/bgr  -hit=main $@

gometalinter --deadline=50s --vendor \
    --cyclo-over=50 --dupl-threshold=100 \
    --exclude=".*should not use dot imports \(golint\)$" \
    --disable-all \
    --enable=vet \
    --enable=deadcode \
    --enable=golint \
    --enable=vetshadow \
    --enable=gocyclo \
    --enable=gofmt \
    --enable=ineffassign \
    --enable=structcheck \
    --skip=.git \
    --skip=.tool \
    --skip=vendor \
    --tests \
    $@
