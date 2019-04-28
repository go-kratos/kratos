package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "binlogv.go",
        "binlogvr.go",
        "bvcsub.go",
        "cmsvideo.go",
        "comment.go",
        "email.go",
        "ftp.go",
        "notice.go",
        "service.go",
        "user.go",
        "video.go",
    ],
    importpath = "go-common/app/job/bbq/video/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/bbq/video/conf:go_default_library",
        "//app/job/bbq/video/dao:go_default_library",
        "//app/job/bbq/video/model:go_default_library",
        "//app/service/bbq/common:go_default_library",
        "//app/service/bbq/notice-service/api/v1:go_default_library",
        "//app/service/bbq/search/api/grpc/v1:go_default_library",
        "//app/service/bbq/topic/api:go_default_library",
        "//app/service/bbq/video/api/grpc/v1:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//vendor/github.com/ftp-master:go_default_library",
        "//vendor/github.com/robfig/cron:go_default_library",
        "//vendor/gopkg.in/gomail.v2:go_default_library",
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
