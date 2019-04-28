package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "binlog.go",
        "binlog_elec.go",
        "elec_syncdb.go",
        "limited_list.go",
        "repaire.go",
        "service.go",
        "task.go",
        "task_account_biz.go",
        "task_account_user.go",
        "task_bill_daily.go",
        "task_bill_monthly.go",
        "task_shell_recharge.go",
        "tool.go",
    ],
    importpath = "go-common/app/job/main/ugcpay/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/ugcpay/conf:go_default_library",
        "//app/job/main/ugcpay/dao:go_default_library",
        "//app/job/main/ugcpay/model:go_default_library",
        "//app/job/main/ugcpay/service/pay:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/robfig/cron:go_default_library",
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
        "//app/job/main/ugcpay/service/pay:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["service_test.go"],
    embed = [":go_default_library"],
    tags = ["automanaged"],
)
