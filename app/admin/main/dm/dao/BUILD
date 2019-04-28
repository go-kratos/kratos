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
        "credit_test.go",
        "dao_test.go",
        "mask_test.go",
        "mc_filter_test.go",
        "mc_subtitle_test.go",
        "message_test.go",
        "mysql_advance_test.go",
        "mysql_filter_test.go",
        "mysql_mask_test.go",
        "mysql_report_test.go",
        "mysql_subtitle_test.go",
        "mysql_task_test.go",
        "mysql_test.go",
        "mysql_transfer_test.go",
        "search_subtitle_test.go",
        "search_test.go",
        "task_test.go",
        "videoup_test.go",
        "workflow_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/dm/conf:go_default_library",
        "//app/admin/main/dm/model:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "account.go",
        "bfs.go",
        "credit.go",
        "dao.go",
        "mask.go",
        "mc_filter.go",
        "mc_subtitle.go",
        "message.go",
        "mysql.go",
        "mysql_advance.go",
        "mysql_filter.go",
        "mysql_mask.go",
        "mysql_report.go",
        "mysql_subtitle.go",
        "mysql_task.go",
        "mysql_transfer.go",
        "search.go",
        "search_subtitle.go",
        "task.go",
        "videoup.go",
        "workflow.go",
    ],
    importpath = "go-common/app/admin/main/dm/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/dm/conf:go_default_library",
        "//app/admin/main/dm/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/database/bfs:go_default_library",
        "//library/database/elastic:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/queue/databus:go_default_library",
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
        "//app/admin/main/dm/dao/oplog:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
