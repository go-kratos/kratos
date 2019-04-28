package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "app_multiple_databus_test.go",
        "app_single_test.go",
        "config_offset_test.go",
        "dao_test.go",
        "es_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/search/conf:go_default_library",
        "//app/job/main/search/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "app_databus.go",
        "app_multiple.go",
        "app_multiple_databus.go",
        "app_single.go",
        "config_asset.go",
        "config_attr.go",
        "config_business.go",
        "config_offset.go",
        "dao.go",
        "es.go",
        "hbase.go",
        "sms.go",
    ],
    importpath = "go-common/app/job/main/search/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/search/conf:go_default_library",
        "//app/job/main/search/model:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/gopkg.in/olivere/elastic.v5:go_default_library",
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
        "//app/job/main/search/dao/base:all-srcs",
        "//app/job/main/search/dao/business:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
