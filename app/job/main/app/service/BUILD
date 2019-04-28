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
        "contribute_test.go",
        "monitor_test.go",
        "retry_test.go",
        "service_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/app/conf:go_default_library",
        "//app/job/main/app/model/space:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "account.go",
        "archive.go",
        "contribute.go",
        "flush.go",
        "monitor.go",
        "retry.go",
        "service.go",
        "stat.go",
    ],
    importpath = "go-common/app/job/main/app/service",
    tags = ["automanaged"],
    deps = [
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/job/main/app/conf:go_default_library",
        "//app/job/main/app/dao/monitor:go_default_library",
        "//app/job/main/app/dao/space:go_default_library",
        "//app/job/main/app/dao/view:go_default_library",
        "//app/job/main/app/model:go_default_library",
        "//app/job/main/app/model/space:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/api/gorpc:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/resource/model:go_default_library",
        "//app/service/main/resource/rpc/client:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/sync/errgroup:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/job/main/app/service/feed:all-srcs",
        "//app/job/main/app/service/show:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
