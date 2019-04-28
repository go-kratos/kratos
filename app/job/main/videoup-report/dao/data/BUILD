package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "admin_test.go",
        "dao_test.go",
        "reply_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/videoup-report/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "admin.go",
        "dao.go",
        "profit.go",
        "reply.go",
        "uper.go",
    ],
    importpath = "go-common/app/job/main/videoup-report/dao/data",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/videoup-report/conf:go_default_library",
        "//app/job/main/videoup-report/model/archive:go_default_library",
        "//app/job/main/videoup-report/model/data:go_default_library",
        "//app/job/main/videoup-report/model/manager:go_default_library",
        "//app/job/main/videoup-report/model/monitor:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
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
