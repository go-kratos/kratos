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
        "//app/admin/main/tag/cmd:all-srcs",
        "//app/admin/main/tag/conf:all-srcs",
        "//app/admin/main/tag/dao:all-srcs",
        "//app/admin/main/tag/http:all-srcs",
        "//app/admin/main/tag/model:all-srcs",
        "//app/admin/main/tag/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
