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
        "//app/admin/main/credit/cmd:all-srcs",
        "//app/admin/main/credit/conf:all-srcs",
        "//app/admin/main/credit/dao/account:all-srcs",
        "//app/admin/main/credit/dao/blocked:all-srcs",
        "//app/admin/main/credit/dao/manager:all-srcs",
        "//app/admin/main/credit/dao/msg:all-srcs",
        "//app/admin/main/credit/dao/relation:all-srcs",
        "//app/admin/main/credit/dao/search:all-srcs",
        "//app/admin/main/credit/dao/upload:all-srcs",
        "//app/admin/main/credit/http:all-srcs",
        "//app/admin/main/credit/model:all-srcs",
        "//app/admin/main/credit/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
