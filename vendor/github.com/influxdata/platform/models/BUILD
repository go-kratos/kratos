package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "inline_fnv.go",
        "inline_strconv_parse.go",
        "points.go",
        "rows.go",
        "statistic.go",
        "time.go",
    ],
    importpath = "github.com/influxdata/platform/models",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/influxdata/platform/pkg/escape:go_default_library"],
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
