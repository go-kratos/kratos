package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "admin_test.go",
        "block_status_test.go",
        "business_test.go",
        "captcha_test.go",
        "config_test.go",
        "content_test.go",
        "credit_user_test.go",
        "cursor_test.go",
        "dao_test.go",
        "databus_test.go",
        "emoji_test.go",
        "memcache_test.go",
        "notice_test.go",
        "redis_test.go",
        "reply_test.go",
        "report_test.go",
        "subject_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/reply/conf:go_default_library",
        "//app/interface/main/reply/model/reply:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "admin.go",
        "block_status.go",
        "business.go",
        "captcha.go",
        "config.go",
        "content.go",
        "credit_user.go",
        "cursor.go",
        "dao.go",
        "databus.go",
        "emoji.go",
        "memcache.go",
        "notice.go",
        "redis.go",
        "reply.go",
        "report.go",
        "subject.go",
    ],
    importpath = "go-common/app/interface/main/reply/dao/reply",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/reply/conf:go_default_library",
        "//app/interface/main/reply/model/reply:go_default_library",
        "//app/interface/main/reply/model/xreply:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/queue/databus:go_default_library",
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
