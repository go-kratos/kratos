package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "cache.go",
        "databus.go",
        "handler.go",
        "handler_resource.go",
        "handler_task.go",
        "monitor.go",
        "report.go",
        "service.go",
        "task_job.go",
        "task_log.go",
        "weight.go",
    ],
    importpath = "go-common/app/job/main/aegis/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/aegis/conf:go_default_library",
        "//app/job/main/aegis/dao:go_default_library",
        "//app/job/main/aegis/dao/email:go_default_library",
        "//app/job/main/aegis/dao/monitor:go_default_library",
        "//app/job/main/aegis/model:go_default_library",
        "//app/job/main/aegis/model/monitor:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/up/api/v1:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/queue/databus/databusutil:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//library/xstr:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "monitor_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/job/main/aegis/conf:go_default_library",
        "//app/job/main/aegis/model/monitor:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/golang/mock/gomock:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
