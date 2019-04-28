package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["discovery_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//library/conf/env:go_default_library",
        "//library/ecode:go_default_library",
        "//library/exp/feature:go_default_library",
        "//library/naming:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "discovery.go",
        "shuffle.go",
    ],
    importpath = "go-common/library/naming/discovery",
    tags = ["automanaged"],
    deps = [
        "//library/conf/env:go_default_library",
        "//library/ecode:go_default_library",
        "//library/exp/feature:go_default_library",
        "//library/log:go_default_library",
        "//library/naming:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/netutil:go_default_library",
        "//library/net/netutil/breaker:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
    ],
)

go_test(
    name = "go_default_xtest",
    srcs = ["example_test.go"],
    tags = ["automanaged"],
    deps = [
        "//library/naming:go_default_library",
        "//library/naming/discovery:go_default_library",
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
