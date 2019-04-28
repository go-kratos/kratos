#!/bin/bash

#set -x

files=`.rider/changefiles.sh "*\.go$"`

exitCode=$?
if [[ ${exitCode} -ne 0 ]]; then
    echo ".rider/changefiles.sh fail"
    exit ${exitCode}
fi

pkgs=""
for file in ${files}
do
    pkg="${file%/*}"
    if [ $? -eq 0 ]; then
        if [[ "${pkgs}" = "" ]]; then
            pkgs="${pkg}"
        else
            pkgs="${pkgs}\n${pkg}"
        fi
    fi
done
echo -e "${pkgs}" > .rider/.pkgs
pkgs=`sort .rider/.pkgs|uniq`

echo -e "${pkgs}"