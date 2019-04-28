package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "hook.go",
        "service.go",
    ],
    importpath = "go-common/app/admin/ep/tapd/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/ep/tapd/conf:go_default_library",
        "//app/admin/ep/tapd/dao:go_default_library",
        "//app/admin/ep/tapd/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//vendor/github.com/robfig/cron:go_default_library",
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
