package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["http.go"],
    importpath = "go-common/library/rate/limit/bench/stress/http",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/rate:go_default_library",
        "//library/rate/limit:go_default_library",
        "//library/rate/limit/bench/stress/conf:go_default_library",
        "//library/rate/limit/bench/stress/service:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
