package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["rpc_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/assist/conf:go_default_library",
        "//app/service/main/assist/model/assist:go_default_library",
        "//app/service/main/assist/service:go_default_library",
        "//vendor/github.com/davecgh/go-spew/spew:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "rpc.go",
        "rpc_assist.go",
        "rpc_log.go",
    ],
    importpath = "go-common/app/service/main/assist/rpc/server",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/assist/conf:go_default_library",
        "//app/service/main/assist/model/assist:go_default_library",
        "//app/service/main/assist/service:go_default_library",
        "//library/net/rpc:go_default_library",
        "//library/net/rpc/context:go_default_library",
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
