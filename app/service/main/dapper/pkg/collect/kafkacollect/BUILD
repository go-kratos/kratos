package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["collect_test.go"],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/main/dapper/model:go_default_library",
        "//app/service/main/dapper/pkg/process:go_default_library",
        "//library/log:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = ["collect.go"],
    importpath = "go-common/app/service/main/dapper/pkg/collect/kafkacollect",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/dapper/model:go_default_library",
        "//app/service/main/dapper/pkg/collect:go_default_library",
        "//app/service/main/dapper/pkg/process:go_default_library",
        "//library/log:go_default_library",
        "//library/stat/prom:go_default_library",
        "//vendor/github.com/Shopify/sarama:go_default_library",
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
