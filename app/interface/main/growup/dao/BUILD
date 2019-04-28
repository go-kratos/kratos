package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "account_test.go",
        "activity_test.go",
        "archive_stat_test.go",
        "av_breach_test.go",
        "banner_test.go",
        "bgm_test.go",
        "column_income_test.go",
        "dao_test.go",
        "exchange_test.go",
        "hbase_test.go",
        "income_test.go",
        "notice_test.go",
        "redis_test.go",
        "special_award_test.go",
        "tool_test.go",
        "up_bill_test.go",
        "up_test.go",
        "withdraw_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/growup/conf:go_default_library",
        "//app/interface/main/growup/model:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "account.go",
        "activity.go",
        "archive_stat.go",
        "av_breach.go",
        "banner.go",
        "bgm.go",
        "column_income.go",
        "dao.go",
        "exchange.go",
        "hbase.go",
        "income.go",
        "notice.go",
        "redis.go",
        "special_award.go",
        "tool.go",
        "up.go",
        "up_bill.go",
        "up_year.go",
        "withdraw.go",
    ],
    importpath = "go-common/app/interface/main/growup/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/growup/conf:go_default_library",
        "//app/interface/main/growup/model:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/interface/openplatform/article/rpc/client:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/account/rpc/client:go_default_library",
        "//app/service/main/vip/model:go_default_library",
        "//app/service/main/vip/rpc/client:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/hbase.v2:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/tsuna/gohbase/hrpc:go_default_library",
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
        "//app/interface/main/growup/dao/newbiedao:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
