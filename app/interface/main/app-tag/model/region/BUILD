package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["region.go"],
    importpath = "go-common/app/interface/main/app-tag/model/region",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-tag/model:go_default_library",
        "//app/interface/main/app-tag/model/bangumi:go_default_library",
        "//app/interface/main/app-tag/model/tag:go_default_library",
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
