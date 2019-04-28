package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "account.go",
        "config.go",
        "databus.go",
        "event.go",
        "eventmsg.go",
        "factor.go",
        "factor_group.go",
        "params.go",
        "punishment.go",
        "qcloud.go",
        "risk.go",
        "rpc.go",
        "service.go",
        "spy.go",
        "stat.go",
        "user_event_history.go",
        "user_info.go",
    ],
    importpath = "go-common/app/service/main/spy/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/account/model:go_default_library",
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
