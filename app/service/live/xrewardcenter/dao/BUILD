package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["dao.go"],
    importpath = "go-common/app/service/live/xrewardcenter/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/gift/api/liverpc:go_default_library",
        "//app/service/live/room/api/liverpc:go_default_library",
        "//app/service/live/xrewardcenter/conf:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
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
        "//app/service/live/xrewardcenter/dao/anchorReward:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
