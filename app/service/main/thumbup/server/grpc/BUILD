package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["server.go"],
    importpath = "go-common/app/service/main/thumbup/server/grpc",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/thumbup/api:go_default_library",
        "//app/service/main/thumbup/model:go_default_library",
        "//app/service/main/thumbup/service:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "@io_bazel_rules_go//proto/wkt:empty_go_proto",
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
