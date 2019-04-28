package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["service.go"],
    importpath = "go-common/app/service/live/xanchor/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/xanchor/conf:go_default_library",
        "//app/service/live/xanchor/dao:go_default_library",
        "//app/service/live/xanchor/service/v1:go_default_library",
        "//library/queue/databus:go_default_library",
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
        "//app/service/live/xanchor/service/consumer/v1:all-srcs",
        "//app/service/live/xanchor/service/v1:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
