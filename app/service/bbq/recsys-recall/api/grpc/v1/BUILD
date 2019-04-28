load(
    "@io_bazel_rules_go//proto:def.bzl",
    "go_proto_library",
)
load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

proto_library(
    name = "v1_proto",
    srcs = ["api.proto"],
    tags = ["automanaged"],
    deps = [
        "//app/job/bbq/recall/proto:proto_proto",
        "@com_google_protobuf//:empty_proto",
        "@gogo_special_proto//github.com/gogo/protobuf/gogoproto",
    ],
)

go_proto_library(
    name = "v1_go_proto",
    compilers = ["@io_bazel_rules_go//proto:gogofast_grpc"],
    importpath = "go-common/app/service/bbq/recsys-recall/api/grpc/v1",
    proto = ":v1_proto",
    tags = ["automanaged"],
    deps = [
        "//app/job/bbq/recall/proto:proto_go_proto",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
        "@io_bazel_rules_go//proto/wkt:empty_go_proto",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [],
    embed = [":v1_go_proto"],
    importpath = "go-common/app/service/bbq/recsys-recall/api/grpc/v1",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/bbq/recall/proto:go_default_library",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
        "@io_bazel_rules_go//proto/wkt:empty_go_proto",
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
