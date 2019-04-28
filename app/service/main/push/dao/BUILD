package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "client_test.go",
        "dao_test.go",
        "databus_test.go",
        "memcache_test.go",
        "mysql_test.go",
        "push_test.go",
        "redis_test.go",
        "report_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/push/conf:go_default_library",
        "//app/service/main/push/dao/oppo:go_default_library",
        "//app/service/main/push/model:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "client.go",
        "dao.go",
        "databus.go",
        "mc.cache.go",
        "memcache.go",
        "mysql.go",
        "mysql_report.go",
        "mysql_setting.go",
        "push.go",
        "redis.go",
        "report.go",
    ],
    importpath = "go-common/app/service/main/push/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/push/conf:go_default_library",
        "//app/service/main/push/dao/apns2:go_default_library",
        "//app/service/main/push/dao/fcm:go_default_library",
        "//app/service/main/push/dao/huawei:go_default_library",
        "//app/service/main/push/dao/jpush:go_default_library",
        "//app/service/main/push/dao/mi:go_default_library",
        "//app/service/main/push/dao/oppo:go_default_library",
        "//app/service/main/push/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/errgroup:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/service/main/push/dao/apns2:all-srcs",
        "//app/service/main/push/dao/fcm:all-srcs",
        "//app/service/main/push/dao/huawei:all-srcs",
        "//app/service/main/push/dao/jpush:all-srcs",
        "//app/service/main/push/dao/mi:all-srcs",
        "//app/service/main/push/dao/oppo:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
