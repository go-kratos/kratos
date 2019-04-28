package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "appeal_test.go",
        "blocked_mc_test.go",
        "blocked_redis_test.go",
        "blocked_test.go",
        "case_mc_test.go",
        "case_redis_test.go",
        "case_test.go",
        "dao_test.go",
        "databus_test.go",
        "http_test.go",
        "jury_test.go",
        "manager_test.go",
        "mc_test.go",
        "memcache_test.go",
        "mysql_test.go",
        "opinion_test.go",
        "redis_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/credit/conf:go_default_library",
        "//app/interface/main/credit/model:go_default_library",
        "//vendor/github.com/go-sql-driver/mysql:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "appeal.go",
        "blocked.go",
        "blocked_mc.go",
        "blocked_redis.go",
        "case.go",
        "case_mc.go",
        "case_redis.go",
        "dao.go",
        "databus.go",
        "http.go",
        "jury.go",
        "manager.go",
        "memcache.go",
        "mysql.go",
        "opinion.go",
        "redis.go",
    ],
    importpath = "go-common/app/interface/main/credit/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/credit/conf:go_default_library",
        "//app/interface/main/credit/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/sync/errgroup:go_default_library",
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
