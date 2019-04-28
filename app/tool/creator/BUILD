package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_binary",
)

go_library(
    name = "go_default_library",
    srcs = [
        "gen.go",
        "main.go",
        "parser.go",
        "templete.go",
        "upgrade.go",
        "utils.go",
    ],
    importpath = "go-common/app/tool/creator",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/otokaze/mock/mockgen:go_default_library",
        "//vendor/github.com/otokaze/mock/mockgen/model:go_default_library",
        "//vendor/golang.org/x/tools/imports:go_default_library",
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

go_binary(
    name = "creator",
    embed = [":go_default_library"],
    tags = ["automanaged"],
)
