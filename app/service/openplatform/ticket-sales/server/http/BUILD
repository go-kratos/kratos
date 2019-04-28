load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "distrib.go",
        "http.go",
        "promotion.go",
        "promotion_group.go",
        "promotion_order.go",
        "settle.go",
    ],
    importpath = "go-common/app/service/openplatform/ticket-sales/server/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/openplatform/ticket-sales/api/grpc/v1:go_default_library",
        "//app/service/openplatform/ticket-sales/conf:go_default_library",
        "//app/service/openplatform/ticket-sales/model:go_default_library",
        "//app/service/openplatform/ticket-sales/service:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/binding:go_default_library",
        "//library/net/http/blademaster/render:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = ["http_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/openplatform/ticket-sales/conf:go_default_library",
        "//app/service/openplatform/ticket-sales/service:go_default_library",
        "//library/conf/paladin:go_default_library",
        "//library/net/http/blademaster:go_default_library",
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
