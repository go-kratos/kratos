#!/usr/bin/env bash
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

# This command is used by bazel as the workspace_status_command
# to implement build stamping with git information.

set -o errexit
set -o nounset
set -o pipefail

export KRATOS_ROOT=$(dirname "${BASH_SOURCE}")/..

source "${KRATOS_ROOT}/build/lib/version.sh"
kratos::version::get_version_vars

# Prefix with STABLE_ so that these values are saved to stable-status.txt
# instead of volatile-status.txt.
# Stamped rules will be retriggered by changes to stable-status.txt, but not by
# changes to volatile-status.txt.
# IMPORTANT: the camelCase vars should match the lists in hack/lib/version.sh
# and pkg/version/def.bzl.
cat <<EOF
STABLE_BUILD_GIT_COMMIT ${KRATOS_GIT_COMMIT-}
STABLE_BUILD_SCM_STATUS ${KRATOS_GIT_TREE_STATE-}
STABLE_BUILD_SCM_REVISION ${KRATOS_GIT_VERSION-}
STABLE_BUILD_MAJOR_VERSION ${KRATOS_GIT_MAJOR-}
STABLE_BUILD_MINOR_VERSION ${KRATOS_GIT_MINOR-}
STABLE_BUILD_TIME ${KRATOS_BUILD_TIME-}
STABLE_DOCKER_TAG ${KRATOS_GIT_VERSION/+/_}
gitCommit ${KRATOS_GIT_COMMIT-}
gitTreeState ${KRATOS_GIT_TREE_STATE-}
gitVersion ${KRATOS_GIT_VERSION-}
gitMajor ${KRATOS_GIT_MAJOR-}
gitMinor ${KRATOS_GIT_MINOR-}
buildDate $(date \
  ${SOURCE_DATE_EPOCH:+"--date=@${SOURCE_DATE_EPOCH}"} \
 -u +'%Y-%m-%dT%H:%M:%SZ')
EOF
