load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.go",
        "hbase.go",
        "http.go",
        "memcache.go",
        "redis.go",
    ],
    importpath = "go-common/app/service/main/secure/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/secure/conf:go_default_library",
        "//app/service/main/secure/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/hbase.v2:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/tsuna/gohbase/hrpc:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "dao_test.go",
        "hbase_test.go",
        "http_test.go",
        "memcache_test.go",
        "redis_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/secure/conf:go_default_library",
        "//app/service/main/secure/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/bouk/monkey:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
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
