package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "addit_test.go",
        "archive_oper_test.go",
        "archive_recheck_test.go",
        "archive_test.go",
        "audit_test.go",
        "config_test.go",
        "dao_test.go",
        "flow_design_test.go",
        "task_test.go",
        "type_test.go",
        "video_test.go",
        "video_track_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/videoup-report/conf:go_default_library",
        "//app/job/main/videoup-report/model/archive:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "addit.go",
        "archive.go",
        "archive_delay.go",
        "archive_oper.go",
        "archive_recheck.go",
        "audit.go",
        "config.go",
        "dao.go",
        "flow_design.go",
        "report.go",
        "task.go",
        "task_consumer.go",
        "task_dispatch.go",
        "task_dispatch_done.go",
        "task_dispatch_extend.go",
        "task_json_config.go",
        "task_oper_history.go",
        "track.go",
        "type.go",
        "video.go",
        "video_track.go",
    ],
    importpath = "go-common/app/job/main/videoup-report/dao/archive",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/videoup-report/conf:go_default_library",
        "//app/job/main/videoup-report/model/archive:go_default_library",
        "//app/job/main/videoup-report/model/task:go_default_library",
        "//app/job/main/videoup-report/model/utils:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/dgryski/go-farm:go_default_library",
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
