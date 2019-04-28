package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "sea_con_test.go",
        "service_test.go",
        "sync_mc_test.go",
        "zone_index_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/tv/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "del_cont.go",
        "del_season.go",
        "filter.go",
        "license.go",
        "sea_con.go",
        "search_sug.go",
        "service.go",
        "sync_ep.go",
        "sync_mc.go",
        "sync_retry.go",
        "sync_season.go",
        "zone_index.go",
    ],
    importpath = "go-common/app/job/main/tv/service/pgc",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/tv/conf:go_default_library",
        "//app/job/main/tv/dao/app:go_default_library",
        "//app/job/main/tv/dao/cms:go_default_library",
        "//app/job/main/tv/dao/ftp:go_default_library",
        "//app/job/main/tv/dao/lic:go_default_library",
        "//app/job/main/tv/dao/playurl:go_default_library",
        "//app/job/main/tv/model/common:go_default_library",
        "//app/job/main/tv/model/pgc:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/robfig/cron:go_default_library",
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
