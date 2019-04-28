#!/bin/bash

#set -x

pkgs=`.rider/changepkgs.sh|grep -v ^vendor/`

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

if [ ! -d "${CI_PROJECT_DIR}/../src" ];then
    mkdir ${CI_PROJECT_DIR}/../src
fi
ln -fs ${CI_PROJECT_DIR} ${CI_PROJECT_DIR}/../src
export GOPATH=${CI_PROJECT_DIR}/..
echo "GOPATH: $GOPATH"
cd $GOPATH/src/go-common

exitCode=0
echo -e "\ngometalinter:"

output=`gometalinter --config=.rider/.gometalinter.json ${pkgs}`
exitCode=$?
if [[ "${output}" != "" ]]; then
    exitCode=1
    echo -e "${output}"
fi

exit ${exitCode}

