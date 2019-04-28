package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "common.go",
        "databus.go",
        "email.go",
        "message.go",
        "report.go",
        "resource.go",
        "task.go",
    ],
    importpath = "go-common/app/job/main/aegis/model",
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
    srcs = [
        ":package-srcs",
        "//app/job/main/aegis/model/monitor:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
