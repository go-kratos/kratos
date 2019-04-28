package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "archive_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/videoup/conf:go_default_library",
        "//app/interface/main/videoup/model/archive:go_default_library",
        "//library/ecode:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "allow.go",
        "app.go",
        "check.go",
        "check_staff.go",
        "client.go",
        "creator.go",
        "deal.go",
        "filter.go",
        "geetest.go",
        "infoc.go",
        "log.go",
        "pay.go",
        "pre.go",
        "service.go",
        "subtitle.go",
        "web.go",
    ],
    importpath = "go-common/app/interface/main/videoup/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/creative/model/tag:go_default_library",
        "//app/interface/main/videoup/conf:go_default_library",
        "//app/interface/main/videoup/dao/account:go_default_library",
        "//app/interface/main/videoup/dao/archive:go_default_library",
        "//app/interface/main/videoup/dao/bfs:go_default_library",
        "//app/interface/main/videoup/dao/creative:go_default_library",
        "//app/interface/main/videoup/dao/dynamic:go_default_library",
        "//app/interface/main/videoup/dao/elec:go_default_library",
        "//app/interface/main/videoup/dao/filter:go_default_library",
        "//app/interface/main/videoup/dao/geetest:go_default_library",
        "//app/interface/main/videoup/dao/mission:go_default_library",
        "//app/interface/main/videoup/dao/order:go_default_library",
        "//app/interface/main/videoup/dao/pay:go_default_library",
        "//app/interface/main/videoup/dao/subtitle:go_default_library",
        "//app/interface/main/videoup/dao/tag:go_default_library",
        "//app/interface/main/videoup/model/archive:go_default_library",
        "//app/interface/main/videoup/model/geetest:go_default_library",
        "//app/interface/main/videoup/model/mission:go_default_library",
        "//app/interface/main/videoup/model/porder:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/log/infoc:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/microcosm-cc/bluemonday:go_default_library",
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
