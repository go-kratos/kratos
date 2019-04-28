package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = ["lint.go"],
    importpath = "go-common/app/admin/main/config/pkg/lint",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_xtest",
    srcs = ["lint_test.go"],
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/config/pkg/lint:go_default_library",
        "//app/admin/main/config/pkg/lint/json:go_default_library",
        "//app/admin/main/config/pkg/lint/toml:go_default_library",
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
        "//app/admin/main/config/pkg/lint/json:all-srcs",
        "//app/admin/main/config/pkg/lint/toml:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
