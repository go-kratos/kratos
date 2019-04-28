package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = ["server.go"],
    importpath = "go-common/app/service/openplatform/ticket-sales/server/grpc",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/openplatform/ticket-sales/api/grpc/v1:go_default_library",
        "//app/service/openplatform/ticket-sales/service:go_default_library",
        "//app/service/openplatform/ticket-sales/service/mis:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = ["server_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/openplatform/ticket-sales/api/grpc/type:go_default_library",
        "//app/service/openplatform/ticket-sales/api/grpc/v1:go_default_library",
        "//app/service/openplatform/ticket-sales/conf:go_default_library",
        "//app/service/openplatform/ticket-sales/service:go_default_library",
        "//library/conf/paladin:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
