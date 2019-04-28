package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "av_charge.go",
        "av_charge_statis.go",
        "bgm.go",
        "column_charge.go",
        "dao.go",
        "date_statia.go",
        "up_charge.go",
    ],
    importpath = "go-common/app/job/main/growup/dao/charge",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/growup/conf:go_default_library",
        "//app/job/main/growup/model/charge:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "av_charge_statis_test.go",
        "av_charge_test.go",
        "bgm_test.go",
        "column_charge_test.go",
        "dao_test.go",
        "date_statia_test.go",
        "up_charge_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/growup/conf:go_default_library",
        "//app/job/main/growup/model/charge:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
