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
        "//app/service/main/push/conf:go_default_library",
        "//app/service/main/push/model:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "business.go",
        "callback.go",
        "progress.go",
        "push.go",
        "report.go",
        "service.go",
        "setting.go",
        "task.go",
        "upload.go",
    ],
    importpath = "go-common/app/service/main/push/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/service/main/filter/model/rpc:go_default_library",
        "//app/service/main/filter/rpc/client:go_default_library",
        "//app/service/main/push/conf:go_default_library",
        "//app/service/main/push/dao:go_default_library",
        "//app/service/main/push/dao/apns2:go_default_library",
        "//app/service/main/push/dao/fcm:go_default_library",
        "//app/service/main/push/dao/huawei:go_default_library",
        "//app/service/main/push/dao/jpush:go_default_library",
        "//app/service/main/push/dao/oppo:go_default_library",
        "//app/service/main/push/model:go_default_library",
        "//library/cache:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
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
