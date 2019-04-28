package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "action_test.go",
        "report_test.go",
        "resource_test.go",
        "service_test.go",
        "sub_sort_test.go",
        "sub_test.go",
        "tag_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/tag/conf:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "action.go",
        "channel.go",
        "check.go",
        "report.go",
        "resource.go",
        "service.go",
        "sub.go",
        "sub_sort.go",
        "tag.go",
    ],
    importpath = "go-common/app/service/main/tag/service",
    tags = ["automanaged"],
    deps = [
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/api/gorpc:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/tag/api:go_default_library",
        "//app/service/main/tag/conf:go_default_library",
        "//app/service/main/tag/dao:go_default_library",
        "//app/service/main/tag/model:go_default_library",
        "//library/cache:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
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
