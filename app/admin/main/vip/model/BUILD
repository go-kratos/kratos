package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dialog.go",
        "jointly.go",
        "model.go",
        "month.go",
        "order.go",
        "panel.go",
        "params.go",
        "platform.go",
        "privilege.go",
        "push.go",
        "query.go",
        "reques.go",
        "resouce.go",
        "tips.go",
        "vip.go",
        "welfare.go",
    ],
    importpath = "go-common/app/admin/main/vip/model",
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
