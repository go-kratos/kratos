package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["supervisor_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = ["supervisor.go"],
    importpath = "go-common/library/net/http/blademaster/middleware/supervisor",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/ecode:go_default_library",
        "//library/net/http/blademaster:go_default_library",
    ],
)

go_test(
    name = "go_default_xtest",
    srcs = ["example_test.go"],
    tags = ["automanaged"],
    deps = [
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/supervisor:go_default_library",
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
