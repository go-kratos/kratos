package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "confirm.go",
        "editor.go",
        "input.go",
        "multiline.go",
        "multiselect.go",
        "password.go",
        "select.go",
        "survey.go",
        "transform.go",
        "validate.go",
    ],
    data = ["Gopkg.toml"],
    importpath = "github.com/AlecAivazis/survey",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/kballard/go-shellquote:go_default_library",
        "//vendor/gopkg.in/AlecAivazis/survey.v1/core:go_default_library",
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
