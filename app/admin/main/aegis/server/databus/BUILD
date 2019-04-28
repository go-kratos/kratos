package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "conf.go",
        "resource.go",
    ],
    importpath = "go-common/app/admin/main/aegis/server/databus",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/aegis/model:go_default_library",
        "//app/admin/main/aegis/model/resource:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = ["resource_test.go"],
    embed = [":go_default_library"],
    tags = ["automanaged"],
)
