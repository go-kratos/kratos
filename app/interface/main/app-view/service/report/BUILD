package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = ["report_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-view/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = ["report.go"],
    importpath = "go-common/app/interface/main/app-view/service/report",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/app-view/conf:go_default_library",
        "//app/interface/main/app-view/dao/archive:go_default_library",
        "//app/interface/main/app-view/dao/elec:go_default_library",
        "//app/interface/main/app-view/dao/report:go_default_library",
        "//app/interface/main/app-view/model:go_default_library",
        "//app/interface/main/app-view/model/elec:go_default_library",
        "//app/interface/main/app-view/model/report:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//library/ecode:go_default_library",
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
