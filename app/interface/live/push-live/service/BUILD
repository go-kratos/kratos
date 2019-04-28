package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "common_message_test.go",
        "mids_test.go",
        "push_test.go",
        "service_test.go",
        "start_live_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/interface/live/push-live/conf:go_default_library",
        "//app/interface/live/push-live/dao:go_default_library",
        "//app/interface/live/push-live/model:go_default_library",
        "//library/cache/redis:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "common_message.go",
        "mids.go",
        "push.go",
        "service.go",
        "start_live.go",
    ],
    importpath = "go-common/app/interface/live/push-live/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/live/push-live/conf:go_default_library",
        "//app/interface/live/push-live/dao:go_default_library",
        "//app/interface/live/push-live/model:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/sync/errgroup:go_default_library",
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
