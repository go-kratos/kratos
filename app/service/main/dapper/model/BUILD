package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "response.go",
        "span.go",
        "spanpoint.go",
        "util.go",
    ],
    importpath = "go-common/app/service/main/dapper/model",
    tags = ["manual"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/net/trace/proto:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
        "@io_bazel_rules_go//proto/wkt:duration_go_proto",
        "@io_bazel_rules_go//proto/wkt:timestamp_go_proto",
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
