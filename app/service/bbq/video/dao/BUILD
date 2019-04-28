package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "archive_filter.go",
        "bvc.go",
        "check.go",
        "comment.go",
        "dao.cache.go",
        "dao.go",
        "databus.go",
        "video.go",
        "video_bvc.go",
        "video_repository.go",
        "video_upload_process.go",
    ],
    importpath = "go-common/app/service/bbq/video/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/bbq/app-bbq/model:go_default_library",
        "//app/service/bbq/common/db/bbq:go_default_library",
        "//app/service/bbq/video/api/grpc/v1:go_default_library",
        "//app/service/bbq/video/conf:go_default_library",
        "//app/service/bbq/video/model:go_default_library",
        "//app/service/bbq/video/model/grpc:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/json-iterator/go:go_default_library",
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
