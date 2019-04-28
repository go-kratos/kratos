package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["client_test.go"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//library/net/rpc/liverpc:go_default_library",
        "//library/net/rpc/liverpc/testdata:go_default_library",
        "//library/net/rpc/liverpc/testdata/v1:go_default_library",
        "//library/net/rpc/liverpc/testdata/v2:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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

go_library(
    name = "go_default_library",
    srcs = ["client.go"],
    importpath = "go-common/library/net/rpc/liverpc/test",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
