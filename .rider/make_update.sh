#!/bin/bash

#set -x

if [ ! -d "${CI_PROJECT_DIR}/../src" ];then
  mkdir ${CI_PROJECT_DIR}/../src
fi
ln -fs ${CI_PROJECT_DIR} ${CI_PROJECT_DIR}/../src
export GOPATH=${CI_PROJECT_DIR}/..
echo "GOPATH: $GOPATH"

vendor=${CI_PROJECT_DIR}/vendor
for dir in $(ls $vendor)
do
  if [ -d $vendor/$dir ]; then
    ln -fs $vendor/$dir $GOPATH/src
  fi
done

sleep $[ ( $RANDOM % 60 )  + 1 ]s

cd ${CI_PROJECT_DIR}/../src/go-common
#mkdir -p .git/hooks/pre-commit
make bazel-update
