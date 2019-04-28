package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "service_test.go",
        "weeklyhonor_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/creative/model/weeklyhonor:go_default_library",
        "//app/job/main/creative/conf:go_default_library",
        "//app/service/main/up/api/v1:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "academy.go",
        "databus.go",
        "service.go",
        "task.go",
        "weeklyhonor.go",
    ],
    importpath = "go-common/app/job/main/creative/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/creative/model/weeklyhonor:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/job/main/creative/conf:go_default_library",
        "//app/job/main/creative/dao/academy:go_default_library",
        "//app/job/main/creative/dao/archive:go_default_library",
        "//app/job/main/creative/dao/monitor:go_default_library",
        "//app/job/main/creative/dao/newcomer:go_default_library",
        "//app/job/main/creative/dao/weeklyhonor:go_default_library",
        "//app/job/main/creative/model:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/up/api/v1:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/xstr:go_default_library",
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
