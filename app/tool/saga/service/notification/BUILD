package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["notification.go"],
    importpath = "go-common/app/tool/saga/service/notification",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/tool/saga/dao:go_default_library",
        "//app/tool/saga/model:go_default_library",
        "//app/tool/saga/service/mail:go_default_library",
        "//app/tool/saga/service/wechat:go_default_library",
        "//library/log:go_default_library",
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
