package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "comment.go",
        "contributors.go",
        "mr.go",
        "pipeline.go",
        "service.go",
        "wechat.go",
    ],
    importpath = "go-common/app/tool/saga/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/tool/saga/conf:go_default_library",
        "//app/tool/saga/dao:go_default_library",
        "//app/tool/saga/model:go_default_library",
        "//app/tool/saga/service/command:go_default_library",
        "//app/tool/saga/service/gitlab:go_default_library",
        "//app/tool/saga/service/mail:go_default_library",
        "//app/tool/saga/service/notification:go_default_library",
        "//app/tool/saga/service/wechat:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/robfig/cron:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/tool/saga/service/command:all-srcs",
        "//app/tool/saga/service/gitlab:all-srcs",
        "//app/tool/saga/service/mail:all-srcs",
        "//app/tool/saga/service/notification:all-srcs",
        "//app/tool/saga/service/wechat:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "comment_test.go",
        "contributors_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/tool/saga/conf:go_default_library",
        "//app/tool/saga/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
