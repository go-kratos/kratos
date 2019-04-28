package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["example_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/tool/gengo/cmd/deepcopy-gen/examples/api/v1:go_default_library",
        "//app/tool/gengo/cmd/deepcopy-gen/examples/model:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = ["example.go"],
    importpath = "go-common/app/tool/gengo/cmd/deepcopy-gen/examples/tests",
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
