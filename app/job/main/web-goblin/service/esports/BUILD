package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = ["service.go"],
    importpath = "go-common/app/job/main/web-goblin/service/esports",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/web-goblin/conf:go_default_library",
        "//app/job/main/web-goblin/dao/esports:go_default_library",
        "//app/job/main/web-goblin/model/esports:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/favorite/api/gorpc:go_default_library",
        "//app/service/main/favorite/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
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
