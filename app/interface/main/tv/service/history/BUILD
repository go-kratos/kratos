package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "combine.go",
        "history.go",
        "service.go",
    ],
    importpath = "go-common/app/interface/main/tv/service/history",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/history/model:go_default_library",
        "//app/interface/main/tv/conf:go_default_library",
        "//app/interface/main/tv/dao/archive:go_default_library",
        "//app/interface/main/tv/dao/cms:go_default_library",
        "//app/interface/main/tv/dao/history:go_default_library",
        "//app/interface/main/tv/model:go_default_library",
        "//app/interface/main/tv/model/history:go_default_library",
        "//library/log:go_default_library",
        "//library/sync/errgroup:go_default_library",
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
    srcs = ["service_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/tv/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
