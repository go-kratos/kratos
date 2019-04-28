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
        "//app/admin/main/passport/cmd:all-srcs",
        "//app/admin/main/passport/conf:all-srcs",
        "//app/admin/main/passport/dao:all-srcs",
        "//app/admin/main/passport/http:all-srcs",
        "//app/admin/main/passport/model:all-srcs",
        "//app/admin/main/passport/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
