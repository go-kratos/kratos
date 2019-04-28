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
    name = "types_proto",
    srcs = ["status.proto"],
    tags = ["automanaged"],
    deps = ["@com_google_protobuf//:any_proto"],
)

go_proto_library(
    name = "types_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_proto"],
    importpath = "go-common/library/ecode/internal/types",
    proto = ":types_proto",
    tags = ["automanaged"],
    deps = ["@io_bazel_rules_go//proto/wkt:any_go_proto"],
)

go_library(
    name = "go_default_library",
    srcs = [],
    embed = [":types_go_proto"],
    importpath = "go-common/library/ecode/internal/types",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_golang_protobuf//proto:go_default_library",
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
