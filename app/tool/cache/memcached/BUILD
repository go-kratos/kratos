package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_binary",
    "go_library",
)

go_binary(
    name = "memcached",
    embed = [":go_default_library"],
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "header_template.go",
        "main.go",
        "multi_template.go",
        "none_template.go",
        "single_template.go",
    ],
    importpath = "go-common/app/tool/cache/memcached",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["//app/tool/cache/common:go_default_library"],
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
        "//app/tool/cache/memcached/testdata:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
