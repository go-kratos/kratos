package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "archive_recheck.go",
        "archive_state.go",
        "config.go",
        "flow.go",
        "forbid.go",
        "history.go",
        "mosaic.go",
        "oper.go",
        "param.go",
        "staff.go",
        "stats.go",
        "task.go",
        "task_qa_video.go",
        "video.go",
        "watermark.go",
    ],
    importpath = "go-common/app/admin/main/videoup/model/archive",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/videoup/model/message:go_default_library",
        "//app/admin/main/videoup/model/utils:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "archive_test.go",
        "history_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = ["//vendor/github.com/smartystreets/goconvey/convey:go_default_library"],
)
