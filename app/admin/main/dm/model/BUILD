package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "action.go",
        "advance.go",
        "dm.go",
        "filter.go",
        "monitor.go",
        "report.go",
        "subject.go",
        "subtitle.go",
        "subtitle_subject.go",
        "task.go",
        "transfer.go",
    ],
    importpath = "go-common/app/admin/main/dm/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
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
    srcs = [
        ":package-srcs",
        "//app/admin/main/dm/model/oplog:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
