package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["aqm_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = ["aqm.go"],
    importpath = "go-common/library/net/http/blademaster/middleware/limit/aqm",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/container/queue/aqm:go_default_library",
        "//library/ecode:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/rate:go_default_library",
        "//library/rate/limit:go_default_library",
        "//library/stat/prom:go_default_library",
    ],
)

go_test(
    name = "go_default_xtest",
    srcs = ["example_test.go"],
    tags = ["automanaged"],
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
