package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "blacklist_test.go",
        "cheat_test.go",
        "combinemails_test.go",
        "data_test.go",
        "httpClient_test.go",
        "income_test.go",
        "recompute_test.go",
        "service_test.go",
        "tag_test.go",
        "tagincome_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/growup/conf:go_default_library",
        "//app/job/main/growup/service/ctrl:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "auto.go",
        "blacklist.go",
        "cheat.go",
        "combinemails.go",
        "creative_activity.go",
        "creative_budget.go",
        "data.go",
        "databus.go",
        "format.go",
        "httpClient.go",
        "income.go",
        "lottery.go",
        "recompute.go",
        "sendincome.go",
        "service.go",
        "signedups.go",
        "tag.go",
        "tagincome.go",
        "task_status.go",
        "tax.go",
        "tool.go",
        "topten.go",
        "up.go",
        "up_bill.go",
        "upload.go",
        "video.go",
    ],
    importpath = "go-common/app/job/main/growup/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/growup/conf:go_default_library",
        "//app/job/main/growup/dao:go_default_library",
        "//app/job/main/growup/dao/charge:go_default_library",
        "//app/job/main/growup/dao/dataplatform:go_default_library",
        "//app/job/main/growup/dao/email:go_default_library",
        "//app/job/main/growup/dao/income:go_default_library",
        "//app/job/main/growup/dao/tag:go_default_library",
        "//app/job/main/growup/model:go_default_library",
        "//app/job/main/growup/model/income:go_default_library",
        "//app/job/main/growup/service/ctrl:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
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
        "//app/job/main/growup/service/charge:all-srcs",
        "//app/job/main/growup/service/ctrl:all-srcs",
        "//app/job/main/growup/service/income:all-srcs",
        "//app/job/main/growup/service/tag:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
