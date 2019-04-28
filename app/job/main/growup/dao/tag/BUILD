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
        "archive_charge_test.go",
        "charge_ratio_test.go",
        "dao_test.go",
        "extra_test.go",
        "income_test.go",
        "tag_info_test.go",
        "tag_up_info_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/growup/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "activity.go",
        "archive_charge.go",
        "charge_ratio.go",
        "dao.go",
        "extra.go",
        "income.go",
        "tag_info.go",
        "tag_up_info.go",
        "up_tag_income.go",
    ],
    importpath = "go-common/app/job/main/growup/dao/tag",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/growup/conf:go_default_library",
        "//app/job/main/growup/model:go_default_library",
        "//app/job/main/growup/model/tag:go_default_library",
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
