package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["args_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = ["//app/tool/gengo/types:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "args.go",
        "var.go",
    ],
    importpath = "go-common/app/tool/gengo/args",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/tool/gengo/generator:go_default_library",
        "//app/tool/gengo/namer:go_default_library",
        "//app/tool/gengo/parser:go_default_library",
        "//app/tool/gengo/types:go_default_library",
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
