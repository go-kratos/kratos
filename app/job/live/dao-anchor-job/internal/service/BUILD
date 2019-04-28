package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "cover_test.go",
        "data_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/job/live/dao-anchor-job/internal/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "cover.go",
        "data.go",
        "service.go",
    ],
    importpath = "go-common/app/job/live/dao-anchor-job/internal/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/live/dao-anchor-job/internal/conf:go_default_library",
        "//app/job/live/dao-anchor-job/internal/dao:go_default_library",
        "//library/log:go_default_library",
        "//library/net/trace:go_default_library",
        "//library/sync/errgroup:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
