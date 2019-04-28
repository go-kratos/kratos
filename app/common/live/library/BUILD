package(default_visibility = ["//visibility:public"])

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
        "//app/common/live/library/lancer:all-srcs",
        "//app/common/live/library/lrucache:all-srcs",
        "//app/common/live/library/mengde:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
