package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["service_test.go"],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/creative/conf:go_default_library",
        "//app/interface/main/creative/service:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "h5.go",
        "service.go",
        "state.go",
        "task.go",
    ],
    importpath = "go-common/app/interface/main/creative/service/newcomer",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/creative/conf:go_default_library",
        "//app/interface/main/creative/dao/academy:go_default_library",
        "//app/interface/main/creative/dao/account:go_default_library",
        "//app/interface/main/creative/dao/archive:go_default_library",
        "//app/interface/main/creative/dao/article:go_default_library",
        "//app/interface/main/creative/dao/data:go_default_library",
        "//app/interface/main/creative/dao/medal:go_default_library",
        "//app/interface/main/creative/dao/newcomer:go_default_library",
        "//app/interface/main/creative/dao/order:go_default_library",
        "//app/interface/main/creative/dao/watermark:go_default_library",
        "//app/interface/main/creative/model/newcomer:go_default_library",
        "//app/interface/main/creative/service:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//library/sync/errgroup.v2:go_default_library",
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
