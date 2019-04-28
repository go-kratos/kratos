package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "api.go",
        "dao.go",
        "pub.go",
    ],
    importpath = "go-common/app/job/live/push-search/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/live/push-search/conf:go_default_library",
        "//app/service/live/relation/api/liverpc:go_default_library",
        "//app/service/live/room/api/liverpc:go_default_library",
        "//app/service/live/user/api/liverpc:go_default_library",
        "//library/database/hbase.v2:go_default_library",
        "//library/log:go_default_library",
        "//library/net/rpc/liverpc:go_default_library",
        "//library/queue/databus:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/job/live/push-search/dao/migrate:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
