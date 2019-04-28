package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "cache.go",
        "conf.go",
        "dao.go",
        "rediskey.go",
    ],
    importpath = "go-common/app/job/live/xroom-feed/internal/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/live/xroom-feed/internal/model:go_default_library",
        "//app/service/live/dao-anchor/api/grpc/v1:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/conf/paladin:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
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
