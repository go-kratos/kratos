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
    name = "version_proto",
    srcs = ["version.proto"],
    tags = ["automanaged"],
)

go_proto_library(
    name = "version_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_proto"],
    importpath = "go-common/app/interface/main/app-resource/model/version",
    proto = ":version_proto",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["version.go"],
    embed = [":version_go_proto"],
    importpath = "go-common/app/interface/main/app-resource/model/version",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/time:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
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
