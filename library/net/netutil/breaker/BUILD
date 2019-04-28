package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "breaker_test.go",
        "sre_breaker_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//library/time:go_default_library",
        "//vendor/github.com/stretchr/testify/assert:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "breaker.go",
        "sre_breaker.go",
    ],
    importpath = "go-common/library/net/netutil/breaker",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/stat/summary:go_default_library",
        "//library/time:go_default_library",
    ],
)

go_test(
    name = "go_default_xtest",
    srcs = ["example_test.go"],
    tags = ["automanaged"],
    deps = [
        "//library/net/netutil/breaker:go_default_library",
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
