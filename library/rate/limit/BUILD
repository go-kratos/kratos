package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["limit_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = ["//library/rate:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = ["limit.go"],
    importpath = "go-common/library/rate/limit",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/container/queue/aqm:go_default_library",
        "//library/log:go_default_library",
        "//library/rate:go_default_library",
        "//library/rate/vegas:go_default_library",
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
        "//library/rate/limit/bench/stress/cmd:all-srcs",
        "//library/rate/limit/bench/stress/conf:all-srcs",
        "//library/rate/limit/bench/stress/dao:all-srcs",
        "//library/rate/limit/bench/stress/http:all-srcs",
        "//library/rate/limit/bench/stress/service:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
