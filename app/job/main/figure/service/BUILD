package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "account_test.go",
        "coin_test.go",
        "reply_test.go",
        "service_test.go",
        "spy_test.go",
        "vip_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/figure/conf:go_default_library",
        "//app/job/main/figure/model:go_default_library",
        "//app/job/main/reply/model/reply:go_default_library",
        "//app/service/main/coin/model:go_default_library",
        "//app/service/main/spy/model:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "account.go",
        "blocked.go",
        "coin.go",
        "danmaku.go",
        "pay.go",
        "reply.go",
        "service.go",
        "spy.go",
        "sync.go",
        "vip.go",
    ],
    importpath = "go-common/app/job/main/figure/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/figure/conf:go_default_library",
        "//app/job/main/figure/dao:go_default_library",
        "//app/job/main/figure/model:go_default_library",
        "//app/service/main/coin/model:go_default_library",
        "//app/service/main/spy/model:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
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
