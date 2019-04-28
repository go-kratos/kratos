#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

export KRATOS_ROOT=$(dirname "${BASH_SOURCE}")/..
source "${KRATOS_ROOT}/build/lib/init.sh"

#kratos::util::ensure-gnu-sed

if  ! bazel version |grep $bazel_version >/dev/null ; then
    kratos::log::info "We suggest you to use bazel $bazel_version for building quickly.
    Mac:           brew upgrade bazel 
    Ubuntu:        sudo apt-get upgrade bazel
    Centos/Redhat: sudo yum update bazel 
    Fedore:        sudo dnf update bazel 
    For more information.Please read this document https://docs.bazel.build/versions/master/install.html 
    " >&2
fi

if [ $(uname -s) = "Linux" ]; then 
        kratos::util::ensure-bazel 
fi

if [ $(uname -s) = "Darwin" ]; 
then 
        kratos::util::ensure-homebrew 
        kratos::util::ensure-homebrew-bazel 
fi

