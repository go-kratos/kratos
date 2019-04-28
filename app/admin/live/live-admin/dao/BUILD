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
        "token.go",
    ],
    importpath = "go-common/app/admin/live/live-admin/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/live/live-admin/conf:go_default_library",
        "//app/service/live/av/api/liverpc:go_default_library",
        "//app/service/live/relation/api/liverpc:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/log:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
