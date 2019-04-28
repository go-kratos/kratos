package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "activity_test.go",
        "auto_test.go",
        "av_breach_test.go",
        "bgm_white_list_test.go",
        "blacklist_test.go",
        "budget_test.go",
        "charge_ratio_test.go",
        "cheat_test.go",
        "dao_test.go",
        "data_test.go",
        "email_test.go",
        "hbase_test.go",
        "income_test.go",
        "lottery_test.go",
        "tag_test.go",
        "task_status_test.go",
        "up_bill_test.go",
        "up_category_info_test.go",
        "up_info_video_test.go",
        "up_quality_info_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/growup/conf:go_default_library",
        "//app/job/main/growup/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "activity.go",
        "auto.go",
        "av_breach.go",
        "bgm_white_list.go",
        "blacklist.go",
        "budget.go",
        "charge_ratio.go",
        "cheat.go",
        "dao.go",
        "data.go",
        "email.go",
        "hbase.go",
        "income.go",
        "lottery.go",
        "tag.go",
        "task_status.go",
        "up_bill.go",
        "up_category_info.go",
        "up_info_video.go",
        "up_quality_info.go",
    ],
    importpath = "go-common/app/job/main/growup/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/growup/conf:go_default_library",
        "//app/job/main/growup/model:go_default_library",
        "//library/database/hbase.v2:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/job/main/growup/dao/charge:all-srcs",
        "//app/job/main/growup/dao/dataplatform:all-srcs",
        "//app/job/main/growup/dao/email:all-srcs",
        "//app/job/main/growup/dao/income:all-srcs",
        "//app/job/main/growup/dao/tag:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
