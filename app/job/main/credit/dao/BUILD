package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "auto_case_test.go",
        "blocked_test.go",
        "case_mc_test.go",
        "case_test.go",
        "dao_test.go",
        "depend_test.go",
        "jury_test.go",
        "memcache_test.go",
        "mysql_test.go",
        "redis_test.go",
        "reply_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/credit/conf:go_default_library",
        "//app/job/main/credit/model:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "auto_case.go",
        "blocked.go",
        "case.go",
        "case_mc.go",
        "dao.go",
        "depend.go",
        "jury.go",
        "memcache.go",
        "mysql.go",
        "redis.go",
        "reply.go",
    ],
    importpath = "go-common/app/job/main/credit/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/credit/conf:go_default_library",
        "//app/job/main/credit/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/time:go_default_library",
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
