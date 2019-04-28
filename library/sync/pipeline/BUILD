package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["pipeline_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//library/net/metadata:go_default_library",
        "//library/time:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = ["pipeline.go"],
    importpath = "go-common/library/sync/pipeline",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/net/metadata:go_default_library",
        "//library/time:go_default_library",
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
        "//library/sync/pipeline/fanout:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
