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
    srcs = [
        "api.proto",
        "archive.proto",
        "reply.proto",
        "request.proto",
        "sign_up.proto",
        "up.proto",
    ],
    tags = ["automanaged"],
    deps = [
        "//app/service/main/archive/api:api_proto",
        "@go_googleapis//google/api:annotations_proto",
        "@gogo_special_proto//github.com/gogo/protobuf/gogoproto",
    ],
)

go_proto_library(
    name = "v1_go_proto",
    compilers = ["@io_bazel_rules_go//proto:gogofast_grpc"],
    importpath = "go-common/app/service/main/up/api/v1",
    proto = ":v1_proto",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/archive/api:api_go_proto",
        "//library/time:go_default_library",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
        "@go_googleapis//google/api:annotations_go_proto",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "api.pb.bm.go",
        "client.go",
        "mock.go",
    ],
    embed = [":v1_go_proto"],
    importpath = "go-common/app/service/main/up/api/v1",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/archive/api:go_default_library",
        "//app/tool/protoc-gen-bm/jsonpb:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/golang/mock/gomock:go_default_library",
        "//vendor/google.golang.org/genproto/googleapis/api/annotations:go_default_library",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
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
