package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "account_notify_test.go",
        "dao_test.go",
        "databus_test.go",
        "exp_test.go",
        "http_test.go",
        "member_log_test.go",
        "member_old_test.go",
        "member_test.go",
        "memcache_test.go",
        "realname_test.go",
        "redis_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/member/conf:go_default_library",
        "//app/job/main/member/model:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/bouk/monkey:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "account_notify.go",
        "dao.go",
        "databus.go",
        "exp.go",
        "hbase.go",
        "http.go",
        "member.go",
        "member_log.go",
        "member_old.go",
        "memcache.go",
        "realname.go",
        "redis.go",
    ],
    importpath = "go-common/app/job/main/member/dao",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/member/conf:go_default_library",
        "//app/job/main/member/dao/block:go_default_library",
        "//app/job/main/member/model:go_default_library",
        "//app/service/main/member/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/queue/databus/report:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/job/main/member/dao/block:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
