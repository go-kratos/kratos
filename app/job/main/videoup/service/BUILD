package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "extra_func_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/videoup/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "audit.go",
        "databus.go",
        "delay.go",
        "extra_func.go",
        "history.go",
        "result.go",
        "retry.go",
        "round.go",
        "send.go",
        "service.go",
        "shot.go",
        "stat.go",
        "video.go",
    ],
    importpath = "go-common/app/job/main/videoup/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/videoup/conf:go_default_library",
        "//app/job/main/videoup/dao/activity:go_default_library",
        "//app/job/main/videoup/dao/archive:go_default_library",
        "//app/job/main/videoup/dao/bvc:go_default_library",
        "//app/job/main/videoup/dao/manager:go_default_library",
        "//app/job/main/videoup/dao/message:go_default_library",
        "//app/job/main/videoup/dao/monitor:go_default_library",
        "//app/job/main/videoup/dao/redis:go_default_library",
        "//app/job/main/videoup/model/archive:go_default_library",
        "//app/job/main/videoup/model/manager:go_default_library",
        "//app/job/main/videoup/model/message:go_default_library",
        "//app/job/main/videoup/model/redis:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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
