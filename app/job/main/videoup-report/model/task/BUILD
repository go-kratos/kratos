package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "task_assign.go",
        "task_dispatch.go",
        "task_weight.go",
    ],
    importpath = "go-common/app/job/main/videoup-report/model/task",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["//app/job/main/videoup-report/model/utils:go_default_library"],
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
