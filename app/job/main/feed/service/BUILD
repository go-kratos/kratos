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
    deps = ["//vendor/github.com/smartystreets/goconvey/convey:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "achive.go",
        "service.go",
    ],
    importpath = "go-common/app/job/main/feed/service",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/feed/conf:go_default_library",
        "//app/job/main/feed/dao:go_default_library",
        "//app/job/main/feed/model:go_default_library",
        "//app/service/main/feed/model:go_default_library",
        "//app/service/main/feed/rpc/client:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
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
