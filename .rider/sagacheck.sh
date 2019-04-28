#!/bin/bash

if [ ! -d "${CI_PROJECT_DIR}/../src" ];then
    mkdir ${CI_PROJECT_DIR}/../src
fi
ln -fs ${CI_PROJECT_DIR} ${CI_PROJECT_DIR}/../src
export GOPATH=${CI_PROJECT_DIR}/..
exitCode=0

# CHANGELOG check
echo "====CHANGELOG check:===="
files=`.rider/changefiles.sh "CHANGELOG.md"`
if [[ "${files}" = "" ]]; then
    echo "未发现CHANGELOG.md文件变更，请'添加'或'修改'CHANGELOG.md"
    exit 1
else
    echo -e "变更如下:\n${files}"
fi

# BGR rule
echo -e "\n====Bili golang rule check:===="
diffFiles=`.rider/changefiles.sh`
cd $GOPATH/src/go-common
go build ./app/tool/bgr
./bgr -script=./app/tool/bgr -hit=main -type=file ${diffFiles}
exitCode=$?

exit ${exitCode}
