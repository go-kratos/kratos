load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "app.go",
        "comment.go",
        "feed.go",
        "location.go",
        "notice.go",
        "push.go",
        "report.go",
        "search.go",
        "service.go",
        "share.go",
        "sv.go",
        "topic.go",
        "upload.go",
        "user.go",
    ],
    importpath = "go-common/app/interface/bbq/app-bbq/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/bbq/app-bbq/api/http/v1:go_default_library",
        "//app/interface/bbq/app-bbq/conf:go_default_library",
        "//app/interface/bbq/app-bbq/dao:go_default_library",
        "//app/interface/bbq/app-bbq/model:go_default_library",
        "//app/interface/bbq/app-bbq/model/grpc:go_default_library",
        "//app/service/bbq/common:go_default_library",
        "//app/service/bbq/notice-service/api/v1:go_default_library",
        "//app/service/bbq/recsys/api/grpc/v1:go_default_library",
        "//app/service/bbq/topic/api:go_default_library",
        "//app/service/bbq/user/api:go_default_library",
        "//app/service/bbq/video/api/grpc/v1:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//library/net/trace:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/Dai0522/go-hash/murmur3:go_default_library",
        "//vendor/github.com/json-iterator/go:go_default_library",
        "@io_bazel_rules_go//proto/wkt:empty_go_proto",
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
