package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)

go_binary(
    name = "client",
    embed = [":go_default_library"],
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "go-common/app/service/bbq/recsys/cmd/client",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/bbq/recsys/api/grpc/v1:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/container/pool:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/Dai0522/go-hash/bloomfilter:go_default_library",
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
