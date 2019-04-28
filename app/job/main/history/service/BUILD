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
    deps = [
        "//app/job/main/history/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "delete.go",
        "merge.go",
        "progress.go",
        "service.go",
    ],
    importpath = "go-common/app/job/main/history/service",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/history/model:go_default_library",
        "//app/job/main/history/conf:go_default_library",
        "//app/job/main/history/dao:go_default_library",
        "//app/job/main/history/model:go_default_library",
        "//app/service/main/history/model:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/pipeline:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/golang.org/x/time/rate:go_default_library",
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

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)
