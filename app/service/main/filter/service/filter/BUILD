package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["filter_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/filter/conf:go_default_library",
        "//app/service/main/filter/model:go_default_library",
        "//app/service/main/filter/service/area:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = ["filter.go"],
    importpath = "go-common/app/service/main/filter/service/filter",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/filter/conf:go_default_library",
        "//app/service/main/filter/model:go_default_library",
        "//app/service/main/filter/model/actriearea:go_default_library",
        "//app/service/main/filter/service/area:go_default_library",
        "//app/service/main/filter/service/regexp:go_default_library",
        "//library/log:go_default_library",
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
