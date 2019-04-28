package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "match_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/activity/conf:go_default_library",
        "//app/job/main/activity/model/like:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "action.go",
        "archive.go",
        "bnj.go",
        "kfc.go",
        "like.go",
        "match.go",
        "service.go",
        "subject.go",
    ],
    importpath = "go-common/app/job/main/activity/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/activity/model:go_default_library",
        "//app/interface/main/activity/model/like:go_default_library",
        "//app/interface/main/activity/rpc/client:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/interface/openplatform/article/rpc/client:go_default_library",
        "//app/job/main/activity/conf:go_default_library",
        "//app/job/main/activity/dao/bnj:go_default_library",
        "//app/job/main/activity/dao/dm:go_default_library",
        "//app/job/main/activity/dao/kfc:go_default_library",
        "//app/job/main/activity/dao/like:go_default_library",
        "//app/job/main/activity/model/bnj:go_default_library",
        "//app/job/main/activity/model/kfc:go_default_library",
        "//app/job/main/activity/model/like:go_default_library",
        "//app/job/main/activity/model/match:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/api/gorpc:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/coin/api/gorpc:go_default_library",
        "//app/service/main/coin/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
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
