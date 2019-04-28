package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "common.go",
        "dao.cache.go",
        "dao.go",
        "extension.go",
        "topic.go",
        "topic_video.go",
    ],
    importpath = "go-common/app/service/bbq/topic/internal/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/bbq/topic/api:go_default_library",
        "//app/service/bbq/topic/internal/model:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/conf/paladin:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/sync/errgroup.v2:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "common_test.go",
        "dao.cache_test.go",
        "dao_test.go",
        "extension_test.go",
        "topic_test.go",
        "topic_video_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/bbq/topic/api:go_default_library",
        "//app/service/bbq/topic/internal/model:go_default_library",
        "//library/conf/paladin:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
