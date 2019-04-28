package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "av_charge_statis_test.go",
        "av_charge_test.go",
        "blacklist_test.go",
        "charge_ratio_test.go",
        "run_test.go",
        "service_test.go",
        "up_account_test.go",
        "up_income_date_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/growup/conf:go_default_library",
        "//app/job/main/growup/model/income:go_default_library",
        "//app/job/main/growup/service/ctrl:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "av_charge.go",
        "av_charge_statis.go",
        "av_income.go",
        "av_income_stat.go",
        "bgm.go",
        "bgm_income.go",
        "bgm_income_stat.go",
        "blacklist.go",
        "business_order.go",
        "charge.go",
        "charge_ratio.go",
        "column_charge.go",
        "column_income.go",
        "column_income_stat.go",
        "common.go",
        "income.go",
        "income_date_statis.go",
        "lottery.go",
        "run.go",
        "run_statis.go",
        "service.go",
        "tax.go",
        "tools.go",
        "up_account.go",
        "up_archive_statis.go",
        "up_income.go",
        "up_income_date.go",
        "up_income_stat.go",
        "up_info_video.go",
    ],
    importpath = "go-common/app/job/main/growup/service/income",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/growup/conf:go_default_library",
        "//app/job/main/growup/dao/dataplatform:go_default_library",
        "//app/job/main/growup/dao/email:go_default_library",
        "//app/job/main/growup/dao/income:go_default_library",
        "//app/job/main/growup/model/income:go_default_library",
        "//app/job/main/growup/service:go_default_library",
        "//app/job/main/growup/service/ctrl:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
        "//vendor/golang.org/x/sync/errgroup:go_default_library",
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
