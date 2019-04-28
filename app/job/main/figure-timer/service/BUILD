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
        "//app/job/main/figure-timer/conf:go_default_library",
        "//app/job/main/figure-timer/dao/mock_dao:go_default_library",
        "//app/job/main/figure-timer/model:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/golang/mock/gomock:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "calc.go",
        "fix.go",
        "rank.go",
        "service.go",
    ],
    importpath = "go-common/app/job/main/figure-timer/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/figure-timer/conf:go_default_library",
        "//app/job/main/figure-timer/dao:go_default_library",
        "//app/job/main/figure-timer/model:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/robfig/cron:go_default_library",
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
