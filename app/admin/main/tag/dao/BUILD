package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "business_test.go",
        "channel_category_test.go",
        "channel_group_test.go",
        "channel_mc_test.go",
        "channel_rule_test.go",
        "channel_test.go",
        "dao_test.go",
        "es_test.go",
        "grpc_test.go",
        "http_test.go",
        "limit_test.go",
        "rank_test.go",
        "region_tag_archive_redis_test.go",
        "relation_mc_test.go",
        "relation_redis_test.go",
        "relation_test.go",
        "report_test.go",
        "resource_log_test.go",
        "synonym_test.go",
        "tag_mc_test.go",
        "tag_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/tag/conf:go_default_library",
        "//app/admin/main/tag/model:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "business.go",
        "channel.go",
        "channel_category.go",
        "channel_group.go",
        "channel_mc.go",
        "channel_rule.go",
        "dao.go",
        "es.go",
        "grpc.go",
        "http.go",
        "limit.go",
        "rank.go",
        "region_tag_archive_redis.go",
        "relation.go",
        "relation_mc.go",
        "relation_redis.go",
        "report.go",
        "resource_log.go",
        "synonym.go",
        "tag.go",
        "tag_mc.go",
    ],
    importpath = "go-common/app/admin/main/tag/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/tag/conf:go_default_library",
        "//app/admin/main/tag/model:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/elastic:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
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
