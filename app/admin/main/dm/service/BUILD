package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "advance_test.go",
        "dm_test.go",
        "filter_test.go",
        "mask_test.go",
        "monitor_test.go",
        "report_test.go",
        "service_test.go",
        "subject_test.go",
        "task_test.go",
        "transfer_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/dm/conf:go_default_library",
        "//app/admin/main/dm/model:go_default_library",
        "//app/admin/main/dm/model/oplog:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//vendor/github.com/davecgh/go-spew/spew:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "advance.go",
        "dm.go",
        "filter.go",
        "mask.go",
        "monitor.go",
        "report.go",
        "service.go",
        "shield.go",
        "subject.go",
        "subtitle.go",
        "subtitle_status.go",
        "subtitle_subject.go",
        "task.go",
        "transfer.go",
    ],
    importpath = "go-common/app/admin/main/dm/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/dm/conf:go_default_library",
        "//app/admin/main/dm/dao:go_default_library",
        "//app/admin/main/dm/dao/oplog:go_default_library",
        "//app/admin/main/dm/model:go_default_library",
        "//app/admin/main/dm/model/oplog:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/api/gorpc:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
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
