package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)

go_binary(
    name = "server",
    embed = [":go_default_library"],
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["server.go"],
    importpath = "go-common/library/net/rpc/warden/benchmark/bench/server",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/net/rpc/warden:go_default_library",
        "//library/net/rpc/warden/benchmark/bench/proto:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/prometheus/client_golang/prometheus/promhttp:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
