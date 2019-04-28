package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "associate_test.go",
        "autorenew_test.go",
        "coupon_test.go",
        "frozen_test.go",
        "order_test.go",
        "salary_test.go",
        "service_test.go",
        "sync_test.go",
        "vip_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/vip/conf:go_default_library",
        "//app/job/main/vip/model:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "associate.go",
        "autorenew.go",
        "bcoin.go",
        "check_data.go",
        "coupon.go",
        "frozen.go",
        "message.go",
        "order.go",
        "push.go",
        "salary.go",
        "salary_data.go",
        "service.go",
        "sync.go",
        "vip.go",
    ],
    importpath = "go-common/app/job/main/vip/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/vip/conf:go_default_library",
        "//app/job/main/vip/dao:go_default_library",
        "//app/job/main/vip/model:go_default_library",
        "//app/service/main/coupon/model:go_default_library",
        "//app/service/main/coupon/rpc/client:go_default_library",
        "//app/service/main/vip/api:go_default_library",
        "//app/service/main/vip/model:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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
