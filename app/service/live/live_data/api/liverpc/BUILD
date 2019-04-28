package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["client.go"],
    importpath = "go-common/app/service/live/live_data/api/liverpc",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/live_data/api/liverpc/v0:go_default_library",
        "//app/service/live/live_data/api/liverpc/v1:go_default_library",
        "//library/net/rpc/liverpc:go_default_library",
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
        "//app/service/live/live_data/api/liverpc/v0:all-srcs",
        "//app/service/live/live_data/api/liverpc/v1:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
