package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "admin.go",
        "reply.go",
        "service.go",
    ],
    importpath = "go-common/app/service/main/reply-feed/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/reply-feed/api:go_default_library",
        "//app/service/main/reply-feed/conf:go_default_library",
        "//app/service/main/reply-feed/dao:go_default_library",
        "//app/service/main/reply-feed/model:go_default_library",
        "//library/log:go_default_library",
        "//library/net/netutil:go_default_library",
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
