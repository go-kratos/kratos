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
    name = "proto_proto",
    srcs = ["forward_index.proto"],
    tags = ["automanaged"],
    deps = [
        "//app/job/bbq/recall/proto/quality:quality_proto",
        "@gogo_special_proto//github.com/gogo/protobuf/gogoproto",
    ],
)

go_proto_library(
    name = "proto_go_proto",
    compilers = ["@io_bazel_rules_go//proto:gogofast_proto"],
    importpath = "go-common/app/job/bbq/recall/proto",
    proto = ":proto_proto",
    tags = ["automanaged"],
    deps = [
        "//app/job/bbq/recall/proto/quality:quality_go_proto",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [],
    embed = [":proto_go_proto"],
    importpath = "go-common/app/job/bbq/recall/proto",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/bbq/recall/proto/quality:go_default_library",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/job/bbq/recall/proto/quality:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
