/*
kratos 是go-common的工具链，提供新项目创建，bazel编译功能

kartos build 本目录之下局部编译，根目录全量编译
NAME:
   kratos build - bazel build

USAGE:
   kratos build [arguments...]

EXAMPLE:
   cd app && kratos build ./service/..  admin  interface/.. tool/cache/...
   kratos build

kartos init 新建新项目
USAGE:
   kratos init [command options] [arguments...]

OPTIONS:
   -d value  部门名
   -t value  项目类型(service,interface,admin,job,common)
   -n value  项目名
   -o value  维护人
   --grpc    是否是GRPC

EXAMPLE:
   kratos init -d main -t service -n test -o wangweizhen
*/
package main
