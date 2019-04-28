package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "archive_feed_test.go",
        "archive_test.go",
        "bangumi_feed_test.go",
        "feed_test.go",
        "mock_test.go",
        "service_test.go",
        "unread_count_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/feed/conf:go_default_library",
        "//app/service/main/feed/dao:go_default_library",
        "//app/service/main/feed/model:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/golang/mock/gomock:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "archive_feed.go",
        "article.go",
        "article_feed.go",
        "bangumi_feed.go",
        "feed.go",
        "new_feed.go",
        "service.go",
        "unread_count.go",
    ],
    importpath = "go-common/app/service/main/feed/service",
    tags = ["automanaged"],
    deps = [
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/interface/openplatform/article/rpc/client:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/account/rpc/client:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/api/gorpc:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/feed/conf:go_default_library",
        "//app/service/main/feed/dao:go_default_library",
        "//app/service/main/feed/model:go_default_library",
        "//library/log:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/time:go_default_library",
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
