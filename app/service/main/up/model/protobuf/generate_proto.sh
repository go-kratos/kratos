gopath="$GOPATH/src"
protobuf_path="$gopath/go-common/vendor/github.com/gogo/protobuf"
echo $protobuf_path
ls "$gopath/go-common/vendor/github.com/gogo/protobuf/gogoproto/"
protoc --gofast_out=".." -I"../" -I"$gopath/go-common/vendor/" ../*.proto
