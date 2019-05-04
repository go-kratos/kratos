#!/bin/bash
set -e
if [[ -z $GOPATH ]]; then
    GOPATH=${HOME}/go
fi
BIN_PATH=$( cut -d ':' -f 1 <<< "$GOPATH" )/bin
if [[ ! -z $GOBIN ]]; then
    BIN_PATH=$GOBIN
fi
if [[ ! -z $INSTALL_PATH ]]; then
    BIN_PATH=$INSTALL_PATH
fi
if [[ -f $BIN_PATH/kprotoc ]]; then
    echo "kprotoc alreay install, remove $BIN_PATH/kprotoc first to reinstall."
    exit 1;
fi

ln -s $GOPATH/src/github.com/bilibili/kratos/tool/kprotoc/kprotoc.sh $BIN_PATH/kprotoc
echo "install kprotoc to $BIN_PATH/kprotoc done!"
