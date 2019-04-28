package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "renderer.go",
        "template.go",
        "write.go",
    ],
    importpath = "gopkg.in/AlecAivazis/survey.v1/core",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/mgutz/ansi:go_default_library",
        "//vendor/gopkg.in/AlecAivazis/survey.v1/terminal:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
