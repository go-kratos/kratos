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
        "//app/service/main/assist/conf:go_default_library",
        "//app/service/main/assist/model/assist:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "assist.go",
        "check.go",
        "limit.go",
        "log.go",
        "relation.go",
        "service.go",
    ],
    importpath = "go-common/app/service/main/assist/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/assist/conf:go_default_library",
        "//app/service/main/assist/dao/account:go_default_library",
        "//app/service/main/assist/dao/assist:go_default_library",
        "//app/service/main/assist/dao/message:go_default_library",
        "//app/service/main/assist/model/assist:go_default_library",
        "//app/service/main/assist/model/message:go_default_library",
        "//app/service/main/assist/model/notify:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
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
