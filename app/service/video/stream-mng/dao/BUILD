package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "backup-stream.go",
        "cache.go",
        "dao.cache.go",
        "dao.go",
        "http.go",
        "localcache.go",
        "main-stream.go",
        "mysql.go",
        "notify.go",
        "official-stream.go",
        "redis.go",
        "stream-change-log.go",
        "upstream-summary.go",
    ],
    importpath = "go-common/app/service/video/stream-mng/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/video/stream-mng/common:go_default_library",
        "//app/service/video/stream-mng/conf:go_default_library",
        "//app/service/video/stream-mng/model:go_default_library",
        "//library/cache:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//vendor/github.com/bluele/gcache:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "backup-stream_test.go",
        "cache_test.go",
        "dao.cache_test.go",
        "dao_test.go",
        "http_test.go",
        "localcache_test.go",
        "main-stream_test.go",
        "mysql_test.go",
        "notify_test.go",
        "official-stream_test.go",
        "redis_test.go",
        "stream-change-log_test.go",
        "upstream-summary_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/video/stream-mng/conf:go_default_library",
        "//app/service/video/stream-mng/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
