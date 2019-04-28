package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
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
        "//app/admin/main/activity/cmd:all-srcs",
        "//app/admin/main/activity/conf:all-srcs",
        "//app/admin/main/activity/dao:all-srcs",
        "//app/admin/main/activity/http:all-srcs",
        "//app/admin/main/activity/model:all-srcs",
        "//app/admin/main/activity/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
