package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "errors_test.go",
        "fsm_test.go",
        "relation_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = ["//vendor/github.com/stretchr/testify/assert:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "errors.go",
        "event.go",
        "fsm.go",
        "realtion.go",
        "utils.go",
    ],
    importpath = "go-common/app/job/main/relation/fsm",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
