package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["base_test.go"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//library/conf/env:go_default_library",
        "//library/naming:go_default_library",
        "//library/net/netutil/breaker:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//library/net/rpc/warden/balancer/wrr:go_default_library",
        "//library/net/rpc/warden/proto/testproto:go_default_library",
        "//library/net/rpc/warden/resolver:go_default_library",
        "//library/time:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
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
        "//library/net/rpc/warden/balancer/wrr/test/client:all-srcs",
        "//library/net/rpc/warden/balancer/wrr/test/server:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)

go_library(
    name = "go_default_library",
    srcs = ["base.go"],
    importpath = "go-common/library/net/rpc/warden/balancer/wrr/test",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
