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
        "//app/admin/main/vip/cmd:all-srcs",
        "//app/admin/main/vip/conf:all-srcs",
        "//app/admin/main/vip/dao:all-srcs",
        "//app/admin/main/vip/http:all-srcs",
        "//app/admin/main/vip/model:all-srcs",
        "//app/admin/main/vip/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
