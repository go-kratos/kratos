package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "addit.go",
        "archive.go",
        "databus.go",
        "delay.go",
        "flow.go",
        "oper.go",
        "recheck.go",
        "report.go",
        "task.go",
        "track.go",
        "video.go",
    ],
    importpath = "go-common/app/job/main/videoup-report/model/archive",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["//library/time:go_default_library"],
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
