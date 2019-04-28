package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["server.go"],
    importpath = "go-common/app/service/main/broadcast/server/grpc",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/broadcast/api/grpc/v1:go_default_library",
        "//app/service/main/broadcast/service:go_default_library",
        "//library/conf/paladin:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_google_grpc//encoding/gzip:go_default_library",
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
