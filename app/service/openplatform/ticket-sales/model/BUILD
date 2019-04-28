load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)
load(
    "@io_bazel_rules_go//proto:def.bzl",
    "go_proto_library",
)

package(default_visibility = ["//visibility:public"])

go_library(
    name = "go_default_library",
    srcs = [
        "cache.go",
        "distrib.go",
        "model.go",
        "order.go",
        "pay.go",
        "promo.go",
        "stock.go",
        "ticket.go",
    ],
    importpath = "go-common/app/service/openplatform/ticket-sales/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/common/openplatform/encoding:go_default_library",
        "//app/service/openplatform/ticket-sales/api/grpc/type:go_default_library",
        "//app/service/openplatform/ticket-sales/api/grpc/v1:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
        "@com_github_gogo_protobuf//types:go_default_library",
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
        "//app/service/openplatform/ticket-sales/model/consts:all-srcs",
        "//app/service/openplatform/ticket-sales/model/order_checker/account:all-srcs",
        "//app/service/openplatform/ticket-sales/model/order_checker/item:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
