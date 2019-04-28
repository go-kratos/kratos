package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["dao.go"],
    importpath = "go-common/app/admin/main/macross/dao/oss",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/macross/conf:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/aliyun/aliyun-oss-go-sdk/oss:go_default_library",
    ],
)

filegroup(
    name = "package-srcs",
    srcs = glob(["**"]),
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
)

filegroup(
    name = "all-srcs",
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
