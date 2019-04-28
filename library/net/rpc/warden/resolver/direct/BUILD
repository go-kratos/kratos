package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["direct_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//library/net/netutil/breaker:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//library/net/rpc/warden/proto/testproto:go_default_library",
        "//library/net/rpc/warden/resolver:go_default_library",
        "//library/time:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = ["direct.go"],
    importpath = "go-common/library/net/rpc/warden/resolver/direct",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/conf/env:go_default_library",
        "//library/naming:go_default_library",
        "//library/net/rpc/warden/resolver:go_default_library",
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
