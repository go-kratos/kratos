#!/usr/bin/env bash
set -e

rm -rf ./a
kratos new a
cd ./a/cmd && go build
if [ $? -ne 0 ]
then
  echo "Failed: all"
  exit 1
fi
rm -rf ./b
kratos new b --grpc
cd ./b/cmd && go build
if [ $? -ne 0 ]
then
  echo "Failed: --grpc"
  exit 1
fi
rm -rf ./c
kratos new c --http
cd ./c/cmd && go build
if [ $? -ne 0 ]
then
  echo "Failed: --http"
  exit 1
fi
