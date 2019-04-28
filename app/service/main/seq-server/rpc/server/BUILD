package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["rpc.go"],
    importpath = "go-common/app/service/main/seq-server/rpc/server",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/seq-server/conf:go_default_library",
        "//app/service/main/seq-server/model:go_default_library",
        "//app/service/main/seq-server/service:go_default_library",
        "//library/net/rpc:go_default_library",
        "//library/net/rpc/context:go_default_library",
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
