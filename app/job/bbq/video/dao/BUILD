package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "berserker.go",
        "cmsvideo.go",
        "comment.go",
        "dao.go",
        "email.go",
        "notice.go",
        "user.go",
        "video.go",
    ],
    importpath = "go-common/app/job/bbq/video/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/bbq/video/conf:go_default_library",
        "//app/job/bbq/video/model:go_default_library",
        "//app/service/bbq/notice-service/api/v1:go_default_library",
        "//app/service/bbq/search/api/grpc/v1:go_default_library",
        "//app/service/bbq/sys-msg/api/v1:go_default_library",
        "//app/service/bbq/video/api/grpc/v1:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/json-iterator/go:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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
