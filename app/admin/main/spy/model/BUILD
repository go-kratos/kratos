package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "config.go",
        "event.go",
        "eventhistory.go",
        "eventmsg.go",
        "factor.go",
        "factor_group.go",
        "log.go",
        "page.go",
        "punishment.go",
        "report.go",
        "service.go",
        "setting.go",
        "stat.go",
        "userinfo.go",
    ],
    importpath = "go-common/app/admin/main/spy/model",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = ["//library/time:go_default_library"],
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
