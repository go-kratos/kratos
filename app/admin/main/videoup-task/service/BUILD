package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "search_test.go",
        "service_test.go",
        "video_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/videoup-task/conf:go_default_library",
        "//app/admin/main/videoup-task/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "consumer.go",
        "extra_func.go",
        "log.go",
        "oper.go",
        "qa_video.go",
        "review.go",
        "search.go",
        "service.go",
        "task.go",
        "task_dispatch.go",
        "task_report.go",
        "utils.go",
        "video.go",
        "weight.go",
    ],
    importpath = "go-common/app/admin/main/videoup-task/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/videoup-task/conf:go_default_library",
        "//app/admin/main/videoup-task/dao:go_default_library",
        "//app/admin/main/videoup-task/model:go_default_library",
        "//app/admin/main/videoup/model/archive:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/up/api/v1:go_default_library",
        "//library/database/elastic:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//library/sync/errgroup:go_default_library",
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
