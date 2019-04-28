package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "bnj.go",
        "service.go",
    ],
    importpath = "go-common/app/interface/main/activity/service/bnj",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/activity/conf:go_default_library",
        "//app/interface/main/activity/dao/bnj:go_default_library",
        "//app/interface/main/activity/dao/like:go_default_library",
        "//app/interface/main/activity/model/bnj:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/sync/errgroup.v2:go_default_library",
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
