package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "alarm.go",
        "app.go",
        "apply.go",
        "auth.go",
        "filter.go",
        "message.go",
        "notify.go",
        "notifygroup.go",
        "oldapp.go",
        "oldauth.go",
        "topic.go",
    ],
    importpath = "go-common/app/admin/main/apm/model/databus",
    tags = ["automanaged"],
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
