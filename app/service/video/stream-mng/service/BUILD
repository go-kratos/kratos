package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "backup.go",
        "change-log.go",
        "change-src.go",
        "check-live.go",
        "clear-cache.go",
        "clear-retweet.go",
        "common.go",
        "main-stream.go",
        "notify.go",
        "official.go",
        "screen-shot.go",
        "service.go",
        "stream.go",
        "stream-cut.go",
        "stream-info.go",
        "stream-mask.go",
        "stream-validate.go",
        "upstream-rtmp.go",
        "upstream-summary.go",
    ],
    importpath = "go-common/app/service/video/stream-mng/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/location/api:go_default_library",
        "//app/service/video/stream-mng/common:go_default_library",
        "//app/service/video/stream-mng/conf:go_default_library",
        "//app/service/video/stream-mng/dao:go_default_library",
        "//app/service/video/stream-mng/model:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//library/time:go_default_library",
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
