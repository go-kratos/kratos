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
        "public_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/creative/conf:go_default_library",
        "//library/cache/redis:go_default_library",
        "//vendor/github.com/davecgh/go-spew/spew:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = ["public.go"],
    importpath = "go-common/app/interface/main/creative/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/creative/conf:go_default_library",
        "//app/interface/main/creative/dao/account:go_default_library",
        "//app/interface/main/creative/dao/activity:go_default_library",
        "//app/interface/main/creative/dao/archive:go_default_library",
        "//app/interface/main/creative/dao/article:go_default_library",
        "//app/interface/main/creative/dao/creative:go_default_library",
        "//app/interface/main/creative/dao/pay:go_default_library",
        "//app/interface/main/creative/dao/subtitle:go_default_library",
        "//app/interface/main/creative/dao/tag:go_default_library",
        "//app/interface/main/creative/dao/up:go_default_library",
        "//app/interface/main/creative/model/activity:go_default_library",
        "//app/interface/main/creative/model/archive:go_default_library",
        "//app/interface/main/creative/model/music:go_default_library",
        "//app/interface/main/creative/model/newcomer:go_default_library",
        "//app/interface/main/creative/model/tag:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
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
        "//app/interface/main/creative/service/academy:all-srcs",
        "//app/interface/main/creative/service/account:all-srcs",
        "//app/interface/main/creative/service/ad:all-srcs",
        "//app/interface/main/creative/service/app:all-srcs",
        "//app/interface/main/creative/service/appeal:all-srcs",
        "//app/interface/main/creative/service/archive:all-srcs",
        "//app/interface/main/creative/service/article:all-srcs",
        "//app/interface/main/creative/service/assist:all-srcs",
        "//app/interface/main/creative/service/danmu:all-srcs",
        "//app/interface/main/creative/service/data:all-srcs",
        "//app/interface/main/creative/service/dynamic:all-srcs",
        "//app/interface/main/creative/service/elec:all-srcs",
        "//app/interface/main/creative/service/faq:all-srcs",
        "//app/interface/main/creative/service/feedback:all-srcs",
        "//app/interface/main/creative/service/geetest:all-srcs",
        "//app/interface/main/creative/service/growup:all-srcs",
        "//app/interface/main/creative/service/medal:all-srcs",
        "//app/interface/main/creative/service/music:all-srcs",
        "//app/interface/main/creative/service/newcomer:all-srcs",
        "//app/interface/main/creative/service/operation:all-srcs",
        "//app/interface/main/creative/service/pay:all-srcs",
        "//app/interface/main/creative/service/reply:all-srcs",
        "//app/interface/main/creative/service/resource:all-srcs",
        "//app/interface/main/creative/service/staff:all-srcs",
        "//app/interface/main/creative/service/template:all-srcs",
        "//app/interface/main/creative/service/up:all-srcs",
        "//app/interface/main/creative/service/version:all-srcs",
        "//app/interface/main/creative/service/watermark:all-srcs",
        "//app/interface/main/creative/service/weeklyhonor:all-srcs",
        "//app/interface/main/creative/service/whitelist:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
