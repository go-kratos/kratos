load(
    "@io_bazel_rules_go//proto:def.bzl",
    "go_proto_library",
)

package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["archive.go"],
    embed = [":archive_go_proto"],
    importpath = "go-common/app/interface/main/app-intl/model/player/archive",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["@com_github_golang_protobuf//proto:go_default_library"],
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

proto_library(
    name = "archive_proto",
    srcs = ["model.proto"],
    tags = ["automanaged"],
)

go_proto_library(
    name = "archive_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_proto"],
    importpath = "go-common/app/interface/main/app-intl/model/player/archive",
    proto = ":archive_proto",
    tags = ["automanaged"],
)
