package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "compare.go",
        "options.go",
        "path.go",
        "reporter.go",
        "unsafe_reflect.go",
    ],
    importpath = "github.com/google/go-cmp/cmp",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//vendor/github.com/google/go-cmp/cmp/internal/diff:go_default_library",
        "//vendor/github.com/google/go-cmp/cmp/internal/function:go_default_library",
        "//vendor/github.com/google/go-cmp/cmp/internal/value:go_default_library",
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
        "//vendor/github.com/google/go-cmp/cmp/internal/diff:all-srcs",
        "//vendor/github.com/google/go-cmp/cmp/internal/function:all-srcs",
        "//vendor/github.com/google/go-cmp/cmp/internal/value:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
