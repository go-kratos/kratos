package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "args_test.go",
        "builder_test.go",
        "cond_test.go",
        "flavor_test.go",
        "modifiers_test.go",
        "select_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "args.go",
        "builder.go",
        "cond.go",
        "flavor.go",
        "modifiers.go",
        "select.go",
    ],
    importpath = "go-common/app/service/main/account-recovery/dao/sqlbuilder",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
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
