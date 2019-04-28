package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "conn_test.go",
        "pubsub_test.go",
        "tcp_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/infra/databus/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "conn.go",
        "pubsub.go",
        "tcp.go",
    ],
    importpath = "go-common/app/infra/databus/tcp",
    tags = ["automanaged"],
    deps = [
        "//app/infra/databus/conf:go_default_library",
        "//app/infra/databus/dsn:go_default_library",
        "//app/infra/databus/model:go_default_library",
        "//app/infra/databus/service:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//vendor/github.com/Shopify/sarama:go_default_library",
        "//vendor/github.com/bsm/sarama-cluster:go_default_library",
        "//vendor/github.com/rcrowley/go-metrics:go_default_library",
        "@com_github_gogo_protobuf//proto:go_default_library",
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
