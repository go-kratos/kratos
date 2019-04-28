package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["deepcopy_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = ["//app/tool/gengo/types:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = ["deepcopy.go"],
    importpath = "go-common/app/tool/gengo/cmd/deepcopy-gen/generators",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/tool/gengo/args:go_default_library",
        "//app/tool/gengo/cmd/deepcopy-gen/sets:go_default_library",
        "//app/tool/gengo/generator:go_default_library",
        "//app/tool/gengo/namer:go_default_library",
        "//app/tool/gengo/types:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
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
