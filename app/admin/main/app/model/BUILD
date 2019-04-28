package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["const.go"],
    importpath = "go-common/app/admin/main/app/model",
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
        "//app/admin/main/app/model/aids:all-srcs",
        "//app/admin/main/app/model/audit:all-srcs",
        "//app/admin/main/app/model/bottom:all-srcs",
        "//app/admin/main/app/model/language:all-srcs",
        "//app/admin/main/app/model/notice:all-srcs",
        "//app/admin/main/app/model/wall:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
