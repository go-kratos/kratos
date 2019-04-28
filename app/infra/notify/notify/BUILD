package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "client_test.go",
        "liverpc_client_test.go",
        "notify_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/infra/notify/conf:go_default_library",
        "//app/infra/notify/dao:go_default_library",
        "//app/infra/notify/model:go_default_library",
        "//library/net/rpc/liverpc:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "client.go",
        "liverpc_client.go",
        "pub.go",
        "sub.go",
    ],
    importpath = "go-common/app/infra/notify/notify",
    tags = ["automanaged"],
    deps = [
        "//app/infra/notify/conf:go_default_library",
        "//app/infra/notify/dao:go_default_library",
        "//app/infra/notify/model:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/netutil:go_default_library",
        "//library/net/rpc/liverpc:go_default_library",
        "//library/stat/prom:go_default_library",
        "//vendor/github.com/Shopify/sarama:go_default_library",
        "//vendor/github.com/bsm/sarama-cluster:go_default_library",
        "//vendor/github.com/rcrowley/go-metrics:go_default_library",
        "//vendor/github.com/snluu/uuid:go_default_library",
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
