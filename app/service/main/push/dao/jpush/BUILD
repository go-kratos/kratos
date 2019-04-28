package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["client_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = ["//vendor/github.com/smartystreets/goconvey/convey:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "audience.go",
        "callback.go",
        "client.go",
        "errcode.go",
        "message.go",
        "notice.go",
        "option.go",
        "payload.go",
        "platform.go",
        "report.go",
        "schedule.go",
    ],
    importpath = "go-common/app/service/main/push/dao/jpush",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/stat:go_default_library",
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
