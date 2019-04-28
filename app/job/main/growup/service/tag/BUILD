package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "service_test.go",
        "tag_test.go",
        "tag_up_info_test.go",
        "up_tag_income_test.go",
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
        "activity.go",
        "archive_charge.go",
        "archive_tag.go",
        "category.go",
        "charge_ratio.go",
        "extra.go",
        "income.go",
        "service.go",
        "tag.go",
        "tag_up_info.go",
        "up_tag.go",
        "up_tag_income.go",
    ],
    importpath = "go-common/app/job/main/growup/service/tag",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/growup/conf:go_default_library",
        "//app/job/main/growup/dao/email:go_default_library",
        "//app/job/main/growup/dao/tag:go_default_library",
        "//app/job/main/growup/model/tag:go_default_library",
        "//app/job/main/growup/service:go_default_library",
        "//app/job/main/growup/service/ctrl:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
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
