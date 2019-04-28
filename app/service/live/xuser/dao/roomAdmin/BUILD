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
        "mc.cache_test.go",
        "mc_test.go",
        "mysql_test.go",
        "redis_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/service/live/xuser/conf:go_default_library",
        "//app/service/live/xuser/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.go",
        "mc.cache.go",
        "mc.go",
        "mysql.go",
        "redis.go",
    ],
    importpath = "go-common/app/service/live/xuser/dao/roomAdmin",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/live/banned_service/api/liverpc/v1:go_default_library",
        "//app/service/live/fans_medal/api/liverpc/v2:go_default_library",
        "//app/service/live/room/api/liverpc/v1:go_default_library",
        "//app/service/live/room/api/liverpc/v2:go_default_library",
        "//app/service/live/xuser/api/grpc/v1:go_default_library",
        "//app/service/live/xuser/conf:go_default_library",
        "//app/service/live/xuser/dao:go_default_library",
        "//app/service/live/xuser/model:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/account/rpc/client:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/orm:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/stat/prom:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
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
