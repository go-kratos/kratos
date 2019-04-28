load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_binary",
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
    importpath = "go-common/app/tool/cache",
    tags = ["automanaged"],
    visibility = ["//visibility:private"],
    deps = ["//app/tool/cache/common:go_default_library"],
)

go_binary(
    name = "cache",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
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
        "//app/tool/cache/common:all-srcs",
        "//app/tool/cache/memcached:all-srcs",
        "//app/tool/cache/testdata:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
