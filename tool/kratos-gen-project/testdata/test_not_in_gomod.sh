#!/usr/bin/env bash
set -e

rm -rf /tmp/test-kratos
mkdir /tmp/test-kratos
kratos new a -d /tmp/test-kratos
cd /tmp/test-kratos/a/cmd && go build
if [ $? -ne 0 ]
then
  echo "Failed: all"
  exit 1
fi
kratos new b -d /tmp/test-kratos --grpc
cd /tmp/test-kratos/b/cmd && go build
if [ $? -ne 0 ]
then
  echo "Failed: --grpc"
  exit 1
fi
kratos new c -d /tmp/test-kratos --http
cd /tmp/test-kratos/c/cmd && go build
if [ $? -ne 0 ]
then
  echo "Failed: --http"
  exit 1
fi
