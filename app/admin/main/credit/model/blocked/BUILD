package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "appeal.go",
        "auto_case.go",
        "case.go",
        "config.go",
        "credit_state.go",
        "info.go",
        "jury.go",
        "kpi.go",
        "labour.go",
        "msg.go",
        "notice.go",
        "opinion.go",
        "param.go",
        "publish.go",
        "reason.go",
    ],
    importpath = "go-common/app/admin/main/credit/model/blocked",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/credit/model:go_default_library",
        "//library/log:go_default_library",
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
