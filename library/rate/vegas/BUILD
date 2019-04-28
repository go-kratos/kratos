package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["vegas_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = ["//library/rate:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "sample.go",
        "vegas.go",
    ],
    importpath = "go-common/library/rate/vegas",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["//library/rate:go_default_library"],
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
