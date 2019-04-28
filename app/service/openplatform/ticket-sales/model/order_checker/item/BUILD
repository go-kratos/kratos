package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["item.go"],
    importpath = "go-common/app/service/openplatform/ticket-sales/model/order_checker/item",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/vip/model:go_default_library",
        "//app/service/openplatform/ticket-item/api/grpc/v1:go_default_library",
        "//app/service/openplatform/ticket-sales/api/grpc/v1:go_default_library",
        "//app/service/openplatform/ticket-sales/dao:go_default_library",
        "//app/service/openplatform/ticket-sales/model/consts:go_default_library",
        "//app/service/openplatform/ticket-sales/model/order_checker/account:go_default_library",
        "//library/ecode:go_default_library",
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
