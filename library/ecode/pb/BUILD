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
    srcs = ["ecode.go"],
    embed = [":pb_go_proto"],
    importpath = "go-common/library/ecode/pb",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/ecode:go_default_library",
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

proto_library(
    name = "pb_proto",
    srcs = ["ecode.proto"],
    tags = ["automanaged"],
    deps = ["@com_google_protobuf//:any_proto"],
)

go_proto_library(
    name = "pb_go_proto",
    compilers = ["@io_bazel_rules_go//proto:go_proto"],
    importpath = "go-common/library/ecode/pb",
    proto = ":pb_proto",
    tags = ["automanaged"],
    deps = ["@io_bazel_rules_go//proto/wkt:any_go_proto"],
)
