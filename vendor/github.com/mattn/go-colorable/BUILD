package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "colorable_others.go",
        "noncolorable.go",
    ],
    importpath = "github.com/mattn/go-colorable",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/mattn/go-isatty:go_default_library"],
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
