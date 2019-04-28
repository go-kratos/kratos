package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["stat.go"],
    importpath = "go-common/library/stat",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["//library/stat/prom:go_default_library"],
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
        "//library/stat/counter:all-srcs",
        "//library/stat/prom:all-srcs",
        "//library/stat/statsd:all-srcs",
        "//library/stat/summary:all-srcs",
        "//library/stat/sys/cpu:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
