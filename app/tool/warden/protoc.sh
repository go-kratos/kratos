#!/bin/bash
DEFAULT_PROTOC_GEN="gogofast"
DEFAULT_PROTOC="protoc"
GO_COMMON_DIR_NAME="go-common"
USR_INCLUDE_DIR="/usr/local/include"

function _install_protoc() {
    osname=$(uname -s)
    echo "install protoc ..."
    case $osname in
        "Darwin" )
            brew install protobuf
            ;;
        *)
            echo "unknown operating system, need install protobuf manual see: https://developers.google.com/protocol-buffers"
            exit 1
            ;;
    esac
}

function _install_protoc_gen() {
    local protoc_gen=$1
    case $protoc_gen in
        "gofast" )
            echo "install gofast from github.com/gogo/protobuf/protoc-gen-gofast"
            go get github.com/gogo/protobuf/protoc-gen-gofast
            ;;
        "gogofast" )
            echo "install gogofast from github.com/gogo/protobuf/protoc-gen-gogofast"
            go get github.com/gogo/protobuf/protoc-gen-gogofast
            ;;
        "gogo" )
            echo "install gogo from github.com/gogo/protobuf/protoc-gen-gogo"
            go get github.com/gogo/protobuf/protoc-gen-gogo
            ;;
        "go" )
            echo "install protoc-gen-go from github.com/golang/protobuf"
            go get github.com/golang/protobuf/{proto,protoc-gen-go}
            ;;
        *)
            echo "can't install protoc-gen-${protoc_gen} automatic !"
            exit 1;
            ;;
    esac
}

function _find_go_common_dir() {
    local go_common_dir_name=$1
    local current_dir=$(pwd)
    while [[ "$(basename $current_dir)" != "$go_common_dir_name" ]]; do
        current_dir=$(dirname $current_dir)
        if [[ "$current_dir" == "/" || "$current_dir" == "." || -z "$current_dir" ]]; then
            return 1
        fi
    done
    echo $current_dir
}

function _fix_pb_file() {
    local target_dir=$1
    echo "fix pb file"
    local pb_files=$(find $target_dir -name "*.pb.go" -type f)
    local pkg_name_esc=$(echo "$target_dir" | sed 's_/_\\/_g')
    for file in $pb_files; do
        echo "fix pb file $file"
        if [[ $(uname -s) == 'Darwin' ]]; then
            sed -i "" -e "s/^import \(.*\) \"app\/\(.*\)\"/import \1 \"go-common\/app\/\2\"/g" $file
        else
            sed -i"" -E "s/^import\s*(.*)\s*\"app\/(.*)\"/import\1\"go-common\/app\/\2\"/g" $file
        fi
    done
}

function _esc_string() {
    echo $(echo "$1" | sed 's_/_\\/_g')
}

function _run_protoc() {
    local proto_dir=$1
    local proto_files=$(find $proto_dir -maxdepth 1 -name "*.proto")
    if [[ -z $proto_files ]]; then
        return
    fi
    local protoc_cmd="$PROTOC -I$PROTO_PATH --${PROTOC_GEN}_out=plugins=grpc:. ${proto_files}"
    echo $protoc_cmd
    $protoc_cmd
}

if [[ -z $PROTOC ]]; then
    PROTOC=${DEFAULT_PROTOC}
    which $PROTOC 
    if [[ "$?" -ne "0" ]]; then
        _install_protoc
    fi
fi
if [[ -z $PROTOC_GEN ]]; then
    PROTOC_GEN=${DEFAULT_PROTOC_GEN}
    which protoc-gen-$PROTOC_GEN
    if [[ "$?" -ne "0" ]]; then
        _install_protoc_gen $PROTOC_GEN
    fi
fi

GO_COMMON_DIR=$(_find_go_common_dir $GO_COMMON_DIR_NAME)
if [[ "$?" != "0" ]]; then
    echo "can't find go-common directoy"
    exit 1
fi

if [[ -z $PROTO_PATH ]]; then
    PROTO_PATH=$GO_COMMON_DIR:$GO_COMMON_DIR/vendor:$USR_INCLUDE_DIR
else
    PROTO_PATH=$PROTO_PATH:$GO_COMMON_DIR:$GO_COMMON_DIR/vendor:$USR_INCLUDE_DIR
fi

if [[ ! -z $1 ]]; then
    cd $1
fi
TARGET_DIR=$(pwd)

GO_COMMON_DIR_ESC=$(_esc_string "$GO_COMMON_DIR/")

TARGET_DIR=${TARGET_DIR//$GO_COMMON_DIR_ESC/}

# switch to go_common
cd $GO_COMMON_DIR

DIRS=$(find $TARGET_DIR -type d)

for dir in $DIRS; do
    echo "run protoc in $dir"
    _run_protoc $dir
done

_fix_pb_file $TARGET_DIR
