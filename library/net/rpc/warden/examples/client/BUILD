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
    srcs = ["client.go"],
    importpath = "go-common/library/net/rpc/warden/examples/client",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/ecode:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//library/net/rpc/warden/proto/testproto:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "@com_github_golang_protobuf//ptypes:go_default_library_gen",
        "@io_bazel_rules_go//proto/wkt:any_go_proto",
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
