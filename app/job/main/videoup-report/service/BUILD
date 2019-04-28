package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "archive_test.go",
        "excitation_test.go",
        "monitor_test.go",
        "reply_test.go",
        "retry_test.go",
        "service_test.go",
        "tag_test.go",
        "ups_test.go",
        "video_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/videoup-report/conf:go_default_library",
        "//app/job/main/videoup-report/model/archive:go_default_library",
        "//app/job/main/videoup-report/model/email:go_default_library",
        "//app/job/main/videoup-report/model/manager:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "archive_result.go",
        "databus.go",
        "email.go",
        "excitation.go",
        "extra_func.go",
        "flow.go",
        "monitor.go",
        "reply.go",
        "retry.go",
        "service.go",
        "tag.go",
        "task.go",
        "task_assign.go",
        "task_dispatch.go",
        "task_weight.go",
        "track.go",
        "ups.go",
        "video.go",
    ],
    importpath = "go-common/app/job/main/videoup-report/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/videoup-report/conf:go_default_library",
        "//app/job/main/videoup-report/dao/archive:go_default_library",
        "//app/job/main/videoup-report/dao/data:go_default_library",
        "//app/job/main/videoup-report/dao/email:go_default_library",
        "//app/job/main/videoup-report/dao/hbase:go_default_library",
        "//app/job/main/videoup-report/dao/manager:go_default_library",
        "//app/job/main/videoup-report/dao/mission:go_default_library",
        "//app/job/main/videoup-report/dao/redis:go_default_library",
        "//app/job/main/videoup-report/dao/tag:go_default_library",
        "//app/job/main/videoup-report/model/archive:go_default_library",
        "//app/job/main/videoup-report/model/email:go_default_library",
        "//app/job/main/videoup-report/model/manager:go_default_library",
        "//app/job/main/videoup-report/model/monitor:go_default_library",
        "//app/job/main/videoup-report/model/task:go_default_library",
        "//app/job/main/videoup-report/model/utils:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/archive/api/gorpc:go_default_library",
        "//app/service/main/up/api/v1:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/xstr:go_default_library",
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
