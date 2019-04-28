package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["season_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "season.go",
        "service.go",
    ],
    importpath = "go-common/app/admin/openplatform/sug/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/openplatform/sug/conf:go_default_library",
        "//app/admin/openplatform/sug/dao:go_default_library",
        "//app/admin/openplatform/sug/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
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
