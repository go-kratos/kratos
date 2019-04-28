package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
)

go_library(
    name = "go_default_library",
    srcs = [
        "item_cf_job.go",
        "service.go",
        "user_area_job.go",
    ],
    importpath = "go-common/app/job/live/recommend-job/internal/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/live/recommend-job/internal/conf:go_default_library",
        "//app/job/live/recommend-job/internal/dao:go_default_library",
        "//app/service/live/recommend/recconst:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/log:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/robfig/cron:go_default_library",
        "//vendor/github.com/siddontang/go/ioutil2:go_default_library",
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
