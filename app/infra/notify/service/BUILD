package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["service_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/infra/notify/conf:go_default_library",
        "//app/infra/notify/model:go_default_library",
        "//library/ecode:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "pub.go",
        "service.go",
    ],
    importpath = "go-common/app/infra/notify/service",
    tags = ["automanaged"],
    deps = [
        "//app/infra/notify/conf:go_default_library",
        "//app/infra/notify/dao:go_default_library",
        "//app/infra/notify/model:go_default_library",
        "//app/infra/notify/notify:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
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
