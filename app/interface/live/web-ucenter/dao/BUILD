package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["dao.go"],
    importpath = "go-common/app/interface/live/web-ucenter/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/live/web-ucenter/conf:go_default_library",
        "//app/service/live/rankdb/api/liverpc:go_default_library",
        "//app/service/live/rc/api/liverpc:go_default_library",
        "//app/service/live/room/api/liverpc:go_default_library",
        "//library/net/rpc/liverpc:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/interface/live/web-ucenter/dao/capsule:all-srcs",
        "//app/interface/live/web-ucenter/dao/history:all-srcs",
        "//app/interface/live/web-ucenter/dao/user:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
