/*
kratos 是Kratos的工具链，提供新项目创建，代码生成等功能

kartos build 本目录之下局部编译，根目录全量编译
NAME:
   kratos build

USAGE:
   kratos build [arguments...]

EXAMPLE:
   cd app && kratos build ./service/..  admin  interface/.. tool/cache/...
   kratos build

kartos init 新建新项目
USAGE:
   kratos init [command options] [arguments...]

OPTIONS:
   -n value  项目名
   -o value  维护人
   --grpc    是否是GRPC

EXAMPLE:
   kratos init -n demo -o kratos
*/
package main
