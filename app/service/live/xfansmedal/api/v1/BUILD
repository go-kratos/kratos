load(
    "@io_bazel_rules_go//proto:def.bzl",
    "go_proto_library",
)

package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

proto_library(
    name = "v1_proto",
    srcs = ["api.proto"],
    tags = ["automanaged"],
)

go_proto_library(
    name = "v1_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_proto"],
    importpath = "go-common/app/service/live/xfansmedal/api/v1",
    proto = ":v1_proto",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["generate.go"],
    embed = [":v1_go_proto"],
    importpath = "go-common/app/service/live/xfansmedal/api/v1",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
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
