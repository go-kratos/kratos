package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["const.go"],
    importpath = "go-common/app/interface/main/app-tag/model",
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
    srcs = [
        ":package-srcs",
        "//app/interface/main/app-tag/model/bangumi:all-srcs",
        "//app/interface/main/app-tag/model/feed:all-srcs",
        "//app/interface/main/app-tag/model/region:all-srcs",
        "//app/interface/main/app-tag/model/tag:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
