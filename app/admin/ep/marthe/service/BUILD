package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "bugly.go",
        "bugly_batch_run.go",
        "mail.go",
        "service.go",
        "tapd.go",
        "user.go",
    ],
    importpath = "go-common/app/admin/ep/marthe/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/ep/marthe/conf:go_default_library",
        "//app/admin/ep/marthe/dao:go_default_library",
        "//app/admin/ep/marthe/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//vendor/github.com/robfig/cron:go_default_library",
        "//vendor/github.com/satori/go.uuid:go_default_library",
        "//vendor/gopkg.in/gomail.v2:go_default_library",
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
