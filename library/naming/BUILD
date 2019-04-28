load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
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
        "//library/naming/discovery:all-srcs",
        "//library/naming/livezk:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)

go_library(
    name = "go_default_library",
    srcs = ["naming.go"],
    importpath = "go-common/library/naming",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
