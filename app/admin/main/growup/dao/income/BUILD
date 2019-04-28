package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "av_breach_test.go",
        "av_income_test.go",
        "bgm_test.go",
        "blacklist_test.go",
        "charge_test.go",
        "dao_test.go",
        "lottery_test.go",
        "up_account_test.go",
        "up_category_info_test.go",
        "up_income_test.go",
        "up_info_video_test.go",
        "up_withdraw_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/growup/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "av_breach.go",
        "av_income.go",
        "bgm.go",
        "blacklist.go",
        "charge.go",
        "dao.go",
        "lottery.go",
        "up_account.go",
        "up_category_info.go",
        "up_income.go",
        "up_info_video.go",
        "up_withdraw.go",
    ],
    importpath = "go-common/app/admin/main/growup/dao/income",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/growup/conf:go_default_library",
        "//app/admin/main/growup/model/income:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
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
