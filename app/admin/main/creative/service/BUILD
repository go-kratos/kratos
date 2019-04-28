package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "service_test.go",
        "upbfs_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/creative/conf:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "academy.go",
        "log.go",
        "material.go",
        "music.go",
        "service.go",
        "task.go",
        "upbfs.go",
        "whitelist.go",
    ],
    importpath = "go-common/app/admin/main/creative/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/creative/conf:go_default_library",
        "//app/admin/main/creative/dao:go_default_library",
        "//app/admin/main/creative/model/academy:go_default_library",
        "//app/admin/main/creative/model/bfs:go_default_library",
        "//app/admin/main/creative/model/logcli:go_default_library",
        "//app/admin/main/creative/model/material:go_default_library",
        "//app/admin/main/creative/model/music:go_default_library",
        "//app/admin/main/creative/model/task:go_default_library",
        "//app/admin/main/creative/model/whitelist:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
        "//vendor/golang.org/x/sync/errgroup:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
