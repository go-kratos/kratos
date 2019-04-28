package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["service.go"],
    importpath = "go-common/app/service/live/fans_medal/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/fans_medal/conf:go_default_library",
        "//app/service/live/fans_medal/dao:go_default_library",
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
        "//app/service/live/fans_medal/service/intimacy:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
