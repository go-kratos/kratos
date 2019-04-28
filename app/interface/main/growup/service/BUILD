package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "income_test.go",
        "notice_test.go",
        "service_test.go",
        "withdraw_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/growup/conf:go_default_library",
        "//app/interface/main/growup/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "activity.go",
        "archive_stat.go",
        "banner.go",
        "bgm.go",
        "column.go",
        "exchange.go",
        "income.go",
        "notice.go",
        "service.go",
        "snow_flake.go",
        "special_award.go",
        "up_bill.go",
        "up_status.go",
        "up_year.go",
        "withdraw.go",
    ],
    importpath = "go-common/app/interface/main/growup/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/growup/conf:go_default_library",
        "//app/interface/main/growup/dao:go_default_library",
        "//app/interface/main/growup/model:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/interface/main/growup/service/newbie:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
