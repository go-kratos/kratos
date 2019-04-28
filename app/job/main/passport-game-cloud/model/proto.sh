#! /bin/sh
# proto.sh
gopath=$GOPATH/src
gogopath=$GOPATH/src/go-common/vendor/github.com/gogo/protobuf
protoc --gofast_out=. --proto_path=$gopath:$gogopath:. *.proto
