package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
    "go_binary",
)

go_test(
    name = "go_default_test",
    srcs = ["bgl_test.go"],
    embed = [":go_default_library"],
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "ast_inspect.go",
        "cache.go",
        "file_filter.go",
        "lint_register.go",
        "main.go",
        "model.go",
        "script_parser.go",
    ],
    importpath = "go-common/app/tool/bgr",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/tool/bgr/log:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/tool/bgr/log:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)

go_binary(
    name = "bgr",
    embed = [":go_default_library"],
    tags = ["automanaged"],
)
