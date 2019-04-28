package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["resource_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "abtest.go",
        "ads.go",
        "archive.go",
        "banner.go",
        "const.go",
        "cpm.go",
        "databus.go",
        "relate.go",
        "resource.go",
        "rpc.go",
        "show.go",
    ],
    importpath = "go-common/app/service/main/resource/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/app-view/model:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
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
