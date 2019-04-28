load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_binary",
)

package(default_visibility = ["//visibility:public"])

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
)

go_library(
    name = "go_default_library",
    srcs = [
        "gdoc.go",
        "swagger.go",
    ],
    importpath = "go-common/app/tool/gdoc",
    tags = ["automanaged"],
)

go_binary(
    name = "gdoc",
    embed = [":go_default_library"],
)
