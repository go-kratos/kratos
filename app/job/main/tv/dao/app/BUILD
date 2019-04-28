package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "cont_del_test.go",
        "cont_sync_test.go",
        "dao_test.go",
        "media_cache_test.go",
        "new_test.go",
        "playurl_test.go",
        "redis_test.go",
        "search_con_test.go",
        "search_sug_test.go",
        "season_del_test.go",
        "season_sync_test.go",
        "sync_mc_test.go",
        "sync_retry_test.go",
        "zone_index_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/tv/conf:go_default_library",
        "//app/job/main/tv/model/common:go_default_library",
        "//app/job/main/tv/model/pgc:go_default_library",
        "//library/database/sql:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "cont_del.go",
        "cont_sync.go",
        "dao.go",
        "media_cache.go",
        "playurl.go",
        "redis.go",
        "search_con.go",
        "search_sug.go",
        "season_del.go",
        "season_sync.go",
        "sync_mc.go",
        "sync_retry.go",
        "zone_index.go",
    ],
    importpath = "go-common/app/job/main/tv/dao/app",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/tv/conf:go_default_library",
        "//app/job/main/tv/model/common:go_default_library",
        "//app/job/main/tv/model/pgc:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/stat/prom:go_default_library",
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
