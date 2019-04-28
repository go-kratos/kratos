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
        "//app/job/main/archive/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "archive_result.go",
        "cache.go",
        "dm.go",
        "retry.go",
        "service.go",
        "video.go",
    ],
    importpath = "go-common/app/job/main/archive/service",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/dm2/model:go_default_library",
        "//app/interface/main/dm2/rpc/client:go_default_library",
        "//app/job/main/archive/conf:go_default_library",
        "//app/job/main/archive/dao/archive:go_default_library",
        "//app/job/main/archive/dao/email:go_default_library",
        "//app/job/main/archive/dao/monitor:go_default_library",
        "//app/job/main/archive/dao/reply:go_default_library",
        "//app/job/main/archive/dao/result:go_default_library",
        "//app/job/main/archive/model/archive:go_default_library",
        "//app/job/main/archive/model/databus:go_default_library",
        "//app/job/main/archive/model/dm:go_default_library",
        "//app/job/main/archive/model/result:go_default_library",
        "//app/job/main/archive/model/retry:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/api/gorpc:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/sql:go_default_library",
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
