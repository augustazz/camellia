#!/usr/bin/env bash

source /etc/profile

file=$1

if [ file'' == '' ]; then
    echo "请输入待处理文件"
    exit
fi

protoc $1 --go_out=.

