package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "config_asset.go",
        "config_attrs.go",
        "config_business.go",
        "config_offset.go",
        "es.go",
        "map_data.go",
        "model.go",
    ],
    importpath = "go-common/app/job/main/search/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
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
