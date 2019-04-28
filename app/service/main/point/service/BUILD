package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "activity_test.go",
        "config_test.go",
        "point_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/point/conf:go_default_library",
        "//app/service/main/point/dao:go_default_library",
        "//app/service/main/point/model:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/bouk/monkey:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "activity.go",
        "config.go",
        "point.go",
        "service.go",
    ],
    importpath = "go-common/app/service/main/point/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/point/conf:go_default_library",
        "//app/service/main/point/dao:go_default_library",
        "//app/service/main/point/model:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
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
