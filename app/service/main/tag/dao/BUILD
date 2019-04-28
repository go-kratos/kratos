package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "action_redis_test.go",
        "action_test.go",
        "channel_mc_test.go",
        "channel_test.go",
        "dao_test.go",
        "limit_test.go",
        "ranking_test.go",
        "report_test.go",
        "resource_log_test.go",
        "resource_mc_test.go",
        "resource_oid_redis_test.go",
        "resource_test.go",
        "sub_redis_test.go",
        "sub_sort_mc_test.go",
        "sub_sort_test.go",
        "sub_test.go",
        "tag_mc_test.go",
        "tag_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/tag/conf:go_default_library",
        "//app/service/main/tag/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "action.go",
        "action_redis.go",
        "channel.go",
        "channel_mc.go",
        "dao.go",
        "limit.go",
        "ranking.go",
        "report.go",
        "resource.go",
        "resource_log.go",
        "resource_mc.go",
        "resource_oid_redis.go",
        "sub.go",
        "sub_redis.go",
        "sub_sort.go",
        "sub_sort_mc.go",
        "tag.go",
        "tag_mc.go",
    ],
    importpath = "go-common/app/service/main/tag/dao",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/tag/conf:go_default_library",
        "//app/service/main/tag/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
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
