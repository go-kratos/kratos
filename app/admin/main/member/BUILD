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
        "//app/admin/main/member/cmd:all-srcs",
        "//app/admin/main/member/conf:all-srcs",
        "//app/admin/main/member/dao:all-srcs",
        "//app/admin/main/member/http:all-srcs",
        "//app/admin/main/member/model:all-srcs",
        "//app/admin/main/member/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
