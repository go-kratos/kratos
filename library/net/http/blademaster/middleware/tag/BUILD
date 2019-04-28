package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = ["tag.go"],
    importpath = "go-common/library/net/http/blademaster/middleware/tag",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["//library/net/http/blademaster:go_default_library"],
)

go_test(
    name = "go_default_xtest",
    srcs = [
        "example_test.go",
        "tag_test.go",
    ],
    tags = ["automanaged"],
    deps = [
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/tag:go_default_library",
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
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
