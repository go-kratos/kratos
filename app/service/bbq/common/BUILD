package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["video_state.go"],
    importpath = "go-common/app/service/bbq/common",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["//library/xstr:go_default_library"],
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
        "//app/service/bbq/common/db/bbq:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
