package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "http_client_test.go",
        "sqltool_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = ["//library/log:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "http_client.go",
        "sqltool.go",
    ],
    importpath = "go-common/app/interface/main/mcn/tool/datacenter",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/log:go_default_library",
        "//library/stat:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
