package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "achievements_test.go",
        "act_mission_test.go",
        "act_test.go",
        "action_test.go",
        "cache_test.go",
        "content_test.go",
        "dao.cache_test.go",
        "dao_test.go",
        "extend_test.go",
        "like_test.go",
        "lottery_test.go",
        "match_redis_test.go",
        "match_test.go",
        "mc.cache_test.go",
        "memcached_test.go",
        "mission_group_test.go",
        "protocol_test.go",
        "redis_test.go",
        "subject_test.go",
        "user_achievements_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/activity/conf:go_default_library",
        "//app/interface/main/activity/model/like:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "achievements.go",
        "act.go",
        "act_mission.go",
        "action.go",
        "cache.go",
        "content.go",
        "dao.cache.go",
        "dao.go",
        "extend.go",
        "like.go",
        "lottery.go",
        "match.go",
        "match_redis.go",
        "mc.cache.go",
        "memcached.go",
        "mission_group.go",
        "protocol.go",
        "redis.go",
        "subject.go",
        "user_achievements.go",
    ],
    importpath = "go-common/app/interface/main/activity/dao/like",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/activity/conf:go_default_library",
        "//app/interface/main/activity/model/like:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/elastic:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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
