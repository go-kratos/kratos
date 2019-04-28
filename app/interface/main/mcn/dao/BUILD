package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
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
        "//app/interface/main/mcn/dao/bfs:all-srcs",
        "//app/interface/main/mcn/dao/cache:all-srcs",
        "//app/interface/main/mcn/dao/datadao:all-srcs",
        "//app/interface/main/mcn/dao/global:all-srcs",
        "//app/interface/main/mcn/dao/mcndao:all-srcs",
        "//app/interface/main/mcn/dao/msg:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
