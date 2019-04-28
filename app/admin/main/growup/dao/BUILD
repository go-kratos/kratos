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
        "authority_test.go",
        "banner_test.go",
        "blacklist_test.go",
        "block_test.go",
        "cheat_test.go",
        "credit_test.go",
        "dao_test.go",
        "expense_info_test.go",
        "notice_test.go",
        "offlineactivity_test.go",
        "special_award_test.go",
        "tag_test.go",
        "trade_test.go",
        "up_tag_income_test.go",
        "up_test.go",
        "upload_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/growup/conf:go_default_library",
        "//app/admin/main/growup/dao/shell:go_default_library",
        "//app/admin/main/growup/model:go_default_library",
        "//app/admin/main/growup/model/offlineactivity:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "activity.go",
        "authority.go",
        "banner.go",
        "blacklist.go",
        "block.go",
        "cheat.go",
        "credit.go",
        "dao.go",
        "expense_info.go",
        "notice.go",
        "offlineactivity.go",
        "special_award.go",
        "tag.go",
        "trade.go",
        "up.go",
        "up_tag_income.go",
        "upload.go",
    ],
    importpath = "go-common/app/admin/main/growup/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/growup/conf:go_default_library",
        "//app/admin/main/growup/dao/shell:go_default_library",
        "//app/admin/main/growup/model:go_default_library",
        "//app/admin/main/growup/model/offlineactivity:go_default_library",
        "//app/admin/main/up/util:go_default_library",
        "//library/database/orm:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/trace:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
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
        "//app/admin/main/growup/dao/income:all-srcs",
        "//app/admin/main/growup/dao/message:all-srcs",
        "//app/admin/main/growup/dao/resource:all-srcs",
        "//app/admin/main/growup/dao/shell:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
