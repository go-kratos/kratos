package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "like_test.go",
        "service_test.go",
        "stats_test.go",
        "user_likes_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/thumbup/conf:go_default_library",
        "//app/service/main/thumbup/model:go_default_library",
        "//library/cache/redis:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "check.go",
        "has_like.go",
        "item_likes.go",
        "like.go",
        "service.go",
        "stats.go",
        "user_likes.go",
    ],
    importpath = "go-common/app/service/main/thumbup/service",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/thumbup/api:go_default_library",
        "//app/service/main/thumbup/conf:go_default_library",
        "//app/service/main/thumbup/dao:go_default_library",
        "//app/service/main/thumbup/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
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
