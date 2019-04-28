load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["user_like_test.go"],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = ["//app/interface/bbq/app-bbq/conf:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "app.go",
        "comment.go",
        "danmu.go",
        "dao.go",
        "feed.go",
        "filter.go",
        "image.go",
        "location.go",
        "notice.go",
        "push.go",
        "redundance.go",
        "report.go",
        "search.go",
        "share.go",
        "sv.go",
        "sv_stat.go",
        "upload.go",
        "user.go",
        "user_like.go",
        "user_relation.go",
        "video.go",
    ],
    importpath = "go-common/app/interface/bbq/app-bbq/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/bbq/app-bbq/api/http/v1:go_default_library",
        "//app/interface/bbq/app-bbq/conf:go_default_library",
        "//app/interface/bbq/app-bbq/model:go_default_library",
        "//app/interface/bbq/app-bbq/model/grpc:go_default_library",
        "//app/service/bbq/common:go_default_library",
        "//app/service/bbq/notice-service/api/v1:go_default_library",
        "//app/service/bbq/recsys/api/grpc/v1:go_default_library",
        "//app/service/bbq/user/api:go_default_library",
        "//app/service/bbq/video-image/api/grpc/v1:go_default_library",
        "//app/service/bbq/video/api/grpc/v1:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/filter/api/grpc/v1:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/json-iterator/go:go_default_library",
        "@org_golang_google_grpc//:go_default_library",
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
