package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "advance_test.go",
        "assist_test.go",
        "dm_test.go",
        "filter_test.go",
        "protect_test.go",
        "report_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/dm/conf:go_default_library",
        "//app/interface/main/dm/model:go_default_library",
        "//library/ecode:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "advance.go",
        "archive.go",
        "assist.go",
        "dm.go",
        "filter.go",
        "protect.go",
        "report.go",
        "service.go",
        "subject.go",
    ],
    importpath = "go-common/app/interface/main/dm/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/dm/conf:go_default_library",
        "//app/interface/main/dm/dao:go_default_library",
        "//app/interface/main/dm/model:go_default_library",
        "//app/interface/main/dm2/model:go_default_library",
        "//app/interface/main/dm2/model/oplog:go_default_library",
        "//app/interface/main/dm2/rpc/client:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/api/gorpc:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/assist/model/assist:go_default_library",
        "//app/service/main/assist/rpc/client:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/zhenjl/cityhash:go_default_library",
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
