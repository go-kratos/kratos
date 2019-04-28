load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "appeal.go",
        "expire.go",
        "notify.go",
        "queue.go",
        "service.go",
        "workflow.go",
    ],
    importpath = "go-common/app/job/main/workflow/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/workflow/model/param:go_default_library",
        "//app/job/main/workflow/conf:go_default_library",
        "//app/job/main/workflow/dao:go_default_library",
        "//app/job/main/workflow/model:go_default_library",
        "//app/service/main/workflow/model:go_default_library",
        "//library/database/elastic:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/sync/errgroup.v2:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = ["service_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/job/main/workflow/conf:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
