package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
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
        "//app/admin/main/aegis/dao/gorm:all-srcs",
        "//app/admin/main/aegis/dao/http:all-srcs",
        "//app/admin/main/aegis/dao/mc:all-srcs",
        "//app/admin/main/aegis/dao/mysql:all-srcs",
        "//app/admin/main/aegis/dao/redis:all-srcs",
        "//app/admin/main/aegis/dao/rpc:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
