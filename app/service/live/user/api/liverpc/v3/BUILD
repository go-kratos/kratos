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
    name = "v3_proto",
    srcs = ["User.proto"],
    tags = ["automanaged"],
    deps = ["@gogo_special_proto//github.com/gogo/protobuf/gogoproto"],
)

go_proto_library(
    name = "v3_go_proto",
    compilers = ["@io_bazel_rules_go//proto:gogofast_grpc"],
    importpath = "go-common/app/service/live/user/api/liverpc/v3",
    proto = ":v3_proto",
    tags = ["automanaged"],
    deps = ["@com_github_gogo_protobuf//gogoproto:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = ["User.liverpc.go"],
    embed = [":v3_go_proto"],
    importpath = "go-common/app/service/live/user/api/liverpc/v3",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/net/rpc/liverpc:go_default_library",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
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
