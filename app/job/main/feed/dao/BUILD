package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["dao_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = ["//vendor/github.com/smartystreets/goconvey/convey:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.go",
        "sms.go",
    ],
    importpath = "go-common/app/job/main/feed/dao",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/feed/conf:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/stat/prom:go_default_library",
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
