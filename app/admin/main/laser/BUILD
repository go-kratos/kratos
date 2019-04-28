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
        "//app/admin/main/laser/cmd:all-srcs",
        "//app/admin/main/laser/conf:all-srcs",
        "//app/admin/main/laser/dao:all-srcs",
        "//app/admin/main/laser/http:all-srcs",
        "//app/admin/main/laser/model:all-srcs",
        "//app/admin/main/laser/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
