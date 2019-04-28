#! /bin/sh
# proto.sh  https://github.com/google/protobuf/releases 下载release包，解压后将include中的文件夹拖到/usr/local/include即可
gopath=$GOPATH/src
gogopath=$GOPATH/src/go-common/vendor/github.com/gogo/protobuf
protoc --gofast_out=. --proto_path=/usr/local/include:$gopath:$gogopath:. vip.proto info.proto card.proto profile.proto usersuit.proto