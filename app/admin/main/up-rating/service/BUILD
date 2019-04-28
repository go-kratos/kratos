package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "authority.go",
        "cache.go",
        "parameter.go",
        "score.go",
        "service.go",
        "statistics.go",
        "transfer.go",
    ],
    importpath = "go-common/app/admin/main/up-rating/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/up-rating/conf:go_default_library",
        "//app/admin/main/up-rating/dao:go_default_library",
        "//app/admin/main/up-rating/dao/global:go_default_library",
        "//app/admin/main/up-rating/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/xstr:go_default_library",
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
