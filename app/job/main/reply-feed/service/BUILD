package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "databus.go",
        "loader.go",
        "reply_set.go",
        "reply_zset.go",
        "service.go",
        "stat.go",
        "statistics.go",
    ],
    importpath = "go-common/app/job/main/reply-feed/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/reply-feed/conf:go_default_library",
        "//app/job/main/reply-feed/dao:go_default_library",
        "//app/job/main/reply-feed/model:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/netutil:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//vendor/github.com/ivpusic/grpool:go_default_library",
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
