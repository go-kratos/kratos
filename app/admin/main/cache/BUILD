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
        "//app/admin/main/cache/cmd:all-srcs",
        "//app/admin/main/cache/conf:all-srcs",
        "//app/admin/main/cache/dao:all-srcs",
        "//app/admin/main/cache/http:all-srcs",
        "//app/admin/main/cache/model:all-srcs",
        "//app/admin/main/cache/script:all-srcs",
        "//app/admin/main/cache/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
