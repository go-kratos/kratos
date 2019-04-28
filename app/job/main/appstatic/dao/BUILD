package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
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
        "//app/job/main/appstatic/dao/caldiff:all-srcs",
        "//app/job/main/appstatic/dao/push:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
