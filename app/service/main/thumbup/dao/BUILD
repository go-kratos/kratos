package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "cache_test.go",
        "dao.cache_test.go",
        "dao_test.go",
        "databus_test.go",
        "item_likes_redis_test.go",
        "memcached_test.go",
        "redis_test.go",
        "tidb_test.go",
        "user_likes_redis_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/thumbup/conf:go_default_library",
        "//app/service/main/thumbup/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "cache.go",
        "dao.cache.go",
        "dao.go",
        "databus.go",
        "item_likes_redis.go",
        "memcached.go",
        "redis.go",
        "tidb.go",
        "user_likes_redis.go",
    ],
    importpath = "go-common/app/service/main/thumbup/dao",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/thumbup/api:go_default_library",
        "//app/service/main/thumbup/conf:go_default_library",
        "//app/service/main/thumbup/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/tidb:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/golang.org/x/sync/singleflight:go_default_library",
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
