package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "cmd.go",
        "creative.go",
        "mcncache.go",
        "reply.go",
        "req.go",
        "tables.go",
        "upcheck.go",
    ],
    importpath = "go-common/app/interface/main/mcn/model/mcnmodel",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/mcn/model:go_default_library",
        "//app/interface/main/mcn/model:go_default_library",
        "//app/interface/main/mcn/model/datamodel:go_default_library",
        "//app/interface/main/mcn/tool/validate:go_default_library",
        "//app/service/main/archive/api:go_default_library",
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
    srcs = ["reply_test.go"],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/mcn/model:go_default_library",
        "//app/interface/main/mcn/model:go_default_library",
    ],
)
