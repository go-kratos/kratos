package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "condrule.go",
        "is_black.go",
        "onlinelist.go",
        "reccache.go",
        "rediskey.go",
        "rule.go",
        "service.go",
        "sort.go",
    ],
    importpath = "go-common/app/job/live/xroom-feed/internal/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/live/xroom-feed/internal/dao:go_default_library",
        "//app/job/live/xroom-feed/internal/model:go_default_library",
        "//app/service/live/dao-anchor/api/grpc/v1:go_default_library",
        "//app/service/live/room/api/liverpc:go_default_library",
        "//app/service/live/room/api/liverpc/v1:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/conf/paladin:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/net/rpc/liverpc:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
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
