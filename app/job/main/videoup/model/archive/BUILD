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
        "archive_result.go",
        "biz_poi.go",
        "biz_vote.go",
        "config.go",
        "cover.go",
        "databus.go",
        "delay.go",
        "forbid.go",
        "oper.go",
        "task.go",
        "type.go",
        "video.go",
    ],
    importpath = "go-common/app/job/main/videoup/model/archive",
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
