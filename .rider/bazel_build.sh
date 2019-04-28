#!/bin/bash

type=${1}

echo "$(date) build..." >> /data/gitlab-runner/build.log

function compileall()
{
    make build-keep-going
}

function compilepart()
{
    pkgs=`.rider/changepkgs.sh`
    exitCode=$?
    if [[ ${exitCode} -ne 0 ]]; then
        echo ".rider/changepkgs.sh fail"
        exit ${exitCode}
    fi
    if [[ "${pkgs}" = "" ]]; then
        echo "no changepkgs"
        exit 0
    fi

    echo -e "change packages:\n${pkgs}\n"

    paths=""
    for pkg in ${pkgs}
    do
        if [[ "${paths}" != "" ]]; then
            paths="${paths} union allpaths(//app/..., //${pkg}:all)"
        else
            paths="allpaths(//app/..., //${pkg}:all)"
        fi
    done

    echo "bazel build..."
    query=`bazel query "${paths}"`
    echo -e "${query}\n"
    bazel build --config=ci --watchfs -k $(echo -e "${query}" | grep -v all-srcs | grep -v package-srcs | grep -v _proto)
    #bazel build $(bazel query "${paths}" |grep -v all-srcs  |grep -v package-srcs)
}

if [[ "${type}" = "part" ]]; then
    compilepart
else
    compileall
fi
