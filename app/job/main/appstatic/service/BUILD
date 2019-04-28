package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "cal_diff.go",
        "push.go",
        "service.go",
        "upbfs.go",
    ],
    importpath = "go-common/app/job/main/appstatic/service",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/appstatic/conf:go_default_library",
        "//app/job/main/appstatic/dao/caldiff:go_default_library",
        "//app/job/main/appstatic/dao/push:go_default_library",
        "//app/job/main/appstatic/model:go_default_library",
        "//library/log:go_default_library",
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
