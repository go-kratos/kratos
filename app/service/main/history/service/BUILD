package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "history.go",
        "service.go",
        "user.go",
    ],
    importpath = "go-common/app/service/main/history/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/history/api/grpc:go_default_library",
        "//app/service/main/history/conf:go_default_library",
        "//app/service/main/history/dao:go_default_library",
        "//app/service/main/history/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/sync/pipeline:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//vendor/golang.org/x/sync/singleflight:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "history_test.go",
        "service_test.go",
        "user_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/history/api/grpc:go_default_library",
        "//app/service/main/history/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
