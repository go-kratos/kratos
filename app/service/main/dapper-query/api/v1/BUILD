load(
    "@io_bazel_rules_go//proto:def.bzl",
    "go_proto_library",
)

package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
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

go_library(
    name = "go_default_library",
    srcs = [
        "api.pb.bm.go",
        "api.pb.go",
        "marshal.go",
    ],
    importpath = "go-common/app/service/main/dapper-query/api/v1",
    tags = ["manual"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/tool/protoc-gen-bm/jsonpb:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
        "@com_github_golang_protobuf//jsonpb:go_default_library",
        "@go_googleapis//google/api:annotations_go_proto",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_x_net//context:go_default_library",
    ],
)
