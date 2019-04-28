load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

package(default_visibility = ["//visibility:public"])

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
)

go_library(
    name = "go_default_library",
    srcs = ["xints.go"],
    importpath = "go-common/app/job/main/relation/model/i64b",
    tags = ["automanaged"],
)

go_test(
    name = "go_default_test",
    srcs = ["xints_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
)
