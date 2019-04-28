package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["service_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = ["//app/job/main/tv/conf:go_default_library"],
)

go_library(
    name = "go_default_library",
    srcs = [
        "service.go",
        "style.go",
    ],
    importpath = "go-common/app/job/main/tv/service/report",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/tv/conf:go_default_library",
        "//app/job/main/tv/dao/report:go_default_library",
        "//app/job/main/tv/model/pgc:go_default_library",
        "//app/job/main/tv/model/report:go_default_library",
        "//library/log:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
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
