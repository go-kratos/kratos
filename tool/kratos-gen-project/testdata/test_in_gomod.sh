#!/usr/bin/env bash
set -e

dir=`pwd`

cd $dir
rm -rf ./a
kratos new a
cd ./a/cmd && go build
if [ $? -ne 0 ]; then
  echo "Failed: all"
  exit 1
else
  rm -rf ../../a
fi

cd $dir
rm -rf ./b
kratos new b --grpc
cd ./b/cmd && go build
if [ $? -ne 0 ];then
  echo "Failed: --grpc"
  exit 1
else
  rm -rf ../../b
fi

cd $dir
rm -rf ./c
kratos new c --http
cd ./c/cmd && go build
if [ $? -ne 0 ]; then
  echo "Failed: --http"
  exit 1
else
  rm -rf ../../c
fi
