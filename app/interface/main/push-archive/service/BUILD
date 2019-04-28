package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "archive_test.go",
        "limit_test.go",
        "service_test.go",
        "setting_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/push-archive/conf:go_default_library",
        "//app/interface/main/push-archive/dao:go_default_library",
        "//app/interface/main/push-archive/model:go_default_library",
        "//app/service/main/push/api/grpc/v1:go_default_library",
        "//library/ecode:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "abtest.go",
        "archive.go",
        "limit.go",
        "relation.go",
        "service.go",
        "setting.go",
        "statistics.go",
    ],
    importpath = "go-common/app/interface/main/push-archive/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/push-archive/conf:go_default_library",
        "//app/interface/main/push-archive/dao:go_default_library",
        "//app/interface/main/push-archive/model:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/account/rpc/client:go_default_library",
        "//app/service/main/push/api/grpc/v1:go_default_library",
        "//library/cache:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
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
