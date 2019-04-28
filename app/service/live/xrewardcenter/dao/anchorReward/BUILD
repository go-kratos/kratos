package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.cache.go",
        "dao.go",
        "mc.cache.go",
        "mc.go",
        "mysql.go",
        "redis.go",
    ],
    importpath = "go-common/app/service/live/xrewardcenter/dao/anchorReward",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/xrewardcenter/conf:go_default_library",
        "//app/service/live/xrewardcenter/model/anchorTask:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/orm:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "dao.cache_test.go",
        "dao_test.go",
        "mc.cache_test.go",
        "mc_test.go",
        "mysql_test.go",
        "redis_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/live/xrewardcenter/conf:go_default_library",
        "//app/service/live/xrewardcenter/model/anchorTask:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
