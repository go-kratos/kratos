package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "active_test.go",
        "retry_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/app-wall/conf:go_default_library",
        "//app/job/main/app-wall/model/offer:go_default_library",
        "//vendor/github.com/bsm/sarama-cluster:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "active.go",
        "retry.go",
        "service.go",
    ],
    importpath = "go-common/app/job/main/app-wall/service/offer",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/app-wall/conf:go_default_library",
        "//app/job/main/app-wall/dao/offer:go_default_library",
        "//app/job/main/app-wall/model:go_default_library",
        "//app/job/main/app-wall/model/offer:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/Shopify/sarama:go_default_library",
        "//vendor/github.com/bsm/sarama-cluster:go_default_library",
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
