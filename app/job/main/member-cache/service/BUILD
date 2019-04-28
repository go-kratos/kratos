package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "account.go",
        "block.go",
        "member.go",
        "service.go",
    ],
    importpath = "go-common/app/job/main/member-cache/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/member-cache/conf:go_default_library",
        "//app/job/main/member-cache/dao:go_default_library",
        "//app/job/main/member-cache/model:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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
