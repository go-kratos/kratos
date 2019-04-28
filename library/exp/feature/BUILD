package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = ["feature_gate.go"],
    importpath = "go-common/library/exp/feature",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["//vendor/github.com/pkg/errors:go_default_library"],
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

go_test(
    name = "go_default_test",
    srcs = ["feature_gate_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
)

go_test(
    name = "go_default_xtest",
    srcs = ["example_test.go"],
    tags = ["automanaged"],
    deps = ["//library/exp/feature:go_default_library"],
)
