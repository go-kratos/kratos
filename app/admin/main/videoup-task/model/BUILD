package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "task_qa_video_test.go",
        "video_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = ["//vendor/github.com/smartystreets/goconvey/convey:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "log.go",
        "oper.go",
        "params.go",
        "review.go",
        "task.go",
        "task_qa.go",
        "task_qa_video.go",
        "time.go",
        "up_group.go",
        "user.go",
        "utils.go",
        "video.go",
        "weight.go",
    ],
    importpath = "go-common/app/admin/main/videoup-task/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
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
