#!/bin/bash
DEFAULT_PROTOC_GEN="gogofast"
DEFAULT_PROTOC="protoc"
KRATOS_DIR_NAME="github.com/bilibili/kratos"
USR_INCLUDE_DIR="/usr/local/include"
GOPATH=$GOPATH
if [[ -z $GOPATH ]]; then
    GOPATH=${HOME}/go
fi

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

function _find_kratos_dir() {
    local kratos_dir_name=$1
    local current_dir="$GOPATH/src/$kratos_dir_name"
    if [[ ! -d $current_dir ]]; then
        go get -u $kratos_dir_name
    fi
    echo $current_dir
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
    local protoc_cmd="$PROTOC -I$PROTO_PATH --${PROTOC_GEN}_out=plugins=grpc:${GOPATH}/src ${proto_files}"
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

KRATOS_DIR=$(_find_kratos_dir $KRATOS_DIR_NAME)
if [[ "$?" != "0" ]]; then
    echo "can't find kratos directoy"
    exit 1
fi

KRATOS_PARENT=$(dirname $KRATOS_DIR)

if [[ -z $PROTO_PATH ]]; then
    PROTO_PATH=$GOPATH/src:$KRATOS_PARENT:$USR_INCLUDE_DIR
else
    PROTO_PATH=$GOPATH/src:$PROTO_PATH:$KRATOS_PARENT:$USR_INCLUDE_DIR
fi

if [[ ! -z $1 ]]; then
    cd $1
fi
TARGET_DIR=$(pwd)

# switch to $GOPATH/src
cd $GOPATH/src
echo "switch workdir to $GOPATH/src"

DIRS=$(find $TARGET_DIR -type d)

for dir in $DIRS; do
    echo "run protoc in $dir"
    _run_protoc $dir
done
