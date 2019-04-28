load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

package(default_visibility = ["//visibility:public"])

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
        "//library/conf/dsn:all-srcs",
        "//library/conf/env:all-srcs",
        "//library/conf/flagvar:all-srcs",
        "//library/conf/paladin:all-srcs",
    ],
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "client.go",
        "client_v2.go",
    ],
    importpath = "go-common/library/conf",
    tags = ["automanaged"],
    deps = [
        "//library/conf/env:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "client_test.go",
        "client_v2_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
)
