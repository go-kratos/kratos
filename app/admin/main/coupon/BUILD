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
        "//app/admin/main/coupon/cmd:all-srcs",
        "//app/admin/main/coupon/conf:all-srcs",
        "//app/admin/main/coupon/dao:all-srcs",
        "//app/admin/main/coupon/http:all-srcs",
        "//app/admin/main/coupon/model:all-srcs",
        "//app/admin/main/coupon/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
