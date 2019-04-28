load(
    "@io_bazel_rules_go//proto:def.bzl",
    "go_proto_library",
)
load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["client.go"],
    embed = [":v1_go_proto"],
    importpath = "go-common/app/service/openplatform/ticket-sales/api/grpc/v1",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/openplatform/ticket-sales/api/grpc/type:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//library/time:go_default_library",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
        "@org_golang_x_net//context:go_default_library",
    ],
)

proto_library(
    name = "v1_proto",
    srcs = [
        "order.proto",
        "pay.proto",
        "promotion.proto",
        "promotion_mis.proto",
        "ticket.proto",
    ],
    tags = ["manual"],
    deps = [
        "//app/service/openplatform/ticket-sales/api/grpc/type:_type_proto",
        "@gogo_special_proto//github.com/gogo/protobuf/gogoproto",
    ],
)

go_proto_library(
    name = "v1_go_proto",
    compilers = ["@io_bazel_rules_go//proto:gogofast_grpc"],
    importpath = "go-common/app/service/openplatform/ticket-sales/api/grpc/v1",
    proto = ":v1_proto",
    tags = ["manual"],
    deps = [
        "//app/service/openplatform/ticket-sales/api/grpc/type:_type_go_proto",
        "//library/time:go_default_library",
        "@com_github_gogo_protobuf//gogoproto:go_default_library",
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
