package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "pool_conf.go",
        "rec_pool.go",
        "service.go",
    ],
    importpath = "go-common/app/service/live/xroom-feed/internal/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/xroom-feed/api:go_default_library",
        "//app/service/live/xroom-feed/internal/dao:go_default_library",
        "//app/service/live/xroom-feed/internal/model:go_default_library",
        "//library/conf/paladin:go_default_library",
        "//library/log:go_default_library",
        "//library/xstr:go_default_library",
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
