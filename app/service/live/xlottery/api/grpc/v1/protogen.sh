#! /bin/sh
# 在环境变量中设置好$GOPATH即可食用
cd $GOPATH/src/go-common/app/service/live/xlottery/api/grpc/v1
$GOPATH/src/go-common/app/tool/warden/protoc.sh
