package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "dynamic_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/dynamic/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "dynamic.go",
        "service.go",
    ],
    importpath = "go-common/app/service/main/dynamic/service",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/api/gorpc:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/dynamic/conf:go_default_library",
        "//app/service/main/dynamic/dao:go_default_library",
        "//library/cache:go_default_library",
        "//library/log:go_default_library",
        "//library/net/metadata:go_default_library",
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
