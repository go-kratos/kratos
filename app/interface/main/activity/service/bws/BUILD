package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "achievement.go",
        "bws.go",
        "check.go",
        "lottery.go",
        "point.go",
    ],
    importpath = "go-common/app/interface/main/activity/service/bws",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/activity/conf:go_default_library",
        "//app/interface/main/activity/dao/bws:go_default_library",
        "//app/interface/main/activity/model/bws:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/usersuit/model:go_default_library",
        "//app/service/main/usersuit/rpc/client:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = ["bws_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/activity/conf:go_default_library",
        "//app/interface/main/activity/model/bws:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
