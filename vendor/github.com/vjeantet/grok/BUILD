package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "graph_test.go",
        "grok_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "graph.go",
        "grok.go",
        "patterns.go",
    ],
    importpath = "github.com/vjeantet/grok",
    tags = ["automanaged"],
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
        "//vendor/github.com/vjeantet/grok/example:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
