load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "distrib.go",
        "order.go",
        "order_detail.go",
        "pay.go",
        "promotion.go",
        "promotion_group.go",
        "promotion_order.go",
        "service.go",
        "stock.go",
        "ticket.go",
    ],
    importpath = "go-common/app/service/openplatform/ticket-sales/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/common/openplatform/random:go_default_library",
        "//app/service/openplatform/ticket-sales/api/grpc/type:go_default_library",
        "//app/service/openplatform/ticket-sales/api/grpc/v1:go_default_library",
        "//app/service/openplatform/ticket-sales/conf:go_default_library",
        "//app/service/openplatform/ticket-sales/dao:go_default_library",
        "//app/service/openplatform/ticket-sales/model:go_default_library",
        "//app/service/openplatform/ticket-sales/model/consts:go_default_library",
        "//app/service/openplatform/ticket-sales/model/order_checker/account:go_default_library",
        "//app/service/openplatform/ticket-sales/model/order_checker/item:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "distrib_test.go",
        "order_test.go",
        "promotion_group_test.go",
        "promotion_order_test.go",
        "promotion_test.go",
        "service_test.go",
        "stock_test.go",
        "ticket_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/openplatform/ticket-sales/api/grpc/type:go_default_library",
        "//app/service/openplatform/ticket-sales/api/grpc/v1:go_default_library",
        "//app/service/openplatform/ticket-sales/conf:go_default_library",
        "//app/service/openplatform/ticket-sales/model:go_default_library",
        "//app/service/openplatform/ticket-sales/model/consts:go_default_library",
        "//library/conf/paladin:go_default_library",
        "//library/ecode:go_default_library",
        "//vendor/github.com/davecgh/go-spew/spew:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
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
        "//app/service/openplatform/ticket-sales/service/mis:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
