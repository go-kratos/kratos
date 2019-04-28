package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["grpc.go"],
    importpath = "go-common/app/service/live/zeus/internal/server/grpc",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/zeus/api/v1:go_default_library",
        "//app/service/live/zeus/internal/conf:go_default_library",
        "//app/service/live/zeus/internal/service/v1:go_default_library",
        "//library/net/rpc/warden:go_default_library",
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
