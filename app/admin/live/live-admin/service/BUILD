package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["service.go"],
    importpath = "go-common/app/admin/live/live-admin/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/live/live-admin/conf:go_default_library",
        "//app/admin/live/live-admin/dao:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/admin/live/live-admin/service/v1:all-srcs",
        "//app/admin/live/live-admin/service/v2:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
