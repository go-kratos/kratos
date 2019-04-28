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
        "//app/admin/main/point/conf:go_default_library",
        "//app/admin/main/point/model:go_default_library",
        "//app/service/main/point/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "point.go",
        "service.go",
    ],
    importpath = "go-common/app/admin/main/point/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/point/conf:go_default_library",
        "//app/admin/main/point/dao:go_default_library",
        "//app/admin/main/point/model:go_default_library",
        "//app/service/main/point/model:go_default_library",
        "//app/service/main/point/rpc/client:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
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
