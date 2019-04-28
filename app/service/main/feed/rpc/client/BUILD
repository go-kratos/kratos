package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["feed_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = ["//app/service/main/feed/model:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = ["feed.go"],
    importpath = "go-common/app/service/main/feed/rpc/client",
    tags = ["automanaged"],
    deps = [
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/service/main/feed/model:go_default_library",
        "//library/net/rpc:go_default_library",
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
