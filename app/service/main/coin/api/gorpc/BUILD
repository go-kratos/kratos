package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["coin_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = ["//app/service/main/coin/model:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "coin.go",
        "mock.go",
    ],
    importpath = "go-common/app/service/main/coin/api/gorpc",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/coin/model:go_default_library",
        "//library/net/rpc:go_default_library",
        "//vendor/github.com/golang/mock/gomock:go_default_library",
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
