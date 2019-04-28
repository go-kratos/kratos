package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "dao_test.go",
        "memcache_test.go",
        "mysql_test.go",
        "redis_test.go",
        "wechat_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/push-strategy/conf:go_default_library",
        "//app/service/main/push/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.go",
        "memcache.go",
        "mysql.go",
        "redis.go",
        "wechat.go",
    ],
    importpath = "go-common/app/service/main/push-strategy/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/push-strategy/conf:go_default_library",
        "//app/service/main/push/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/stat/prom:go_default_library",
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
