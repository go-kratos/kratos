load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)
load(
    "@io_bazel_rules_go//proto:def.bzl",
    "go_proto_library",
)

proto_library(
    name = "grpc_proto",
    srcs = ["bvc.proto"],
    tags = ["automanaged"],
    deps = ["@gogo_special_proto//github.com/gogo/protobuf/gogoproto"],
)

go_proto_library(
    name = "grpc_go_proto",
    compilers = ["@io_bazel_rules_go//proto:gogofast_grpc"],
    importpath = "go-common/app/interface/bbq/app-bbq/model/grpc",
    proto = ":grpc_proto",
    tags = ["automanaged"],
    deps = ["@com_github_gogo_protobuf//gogoproto:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = [],
    embed = [":grpc_go_proto"],
    importpath = "go-common/app/interface/bbq/app-bbq/model/grpc",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
        "@com_github_gogo_protobuf//sortkeys:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_x_net//context:go_default_library",
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
