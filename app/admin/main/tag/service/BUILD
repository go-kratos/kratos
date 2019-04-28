package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "business_test.go",
        "channel_test.go",
        "limit_test.go",
        "rank_test.go",
        "relation_test.go",
        "report_test.go",
        "resource_log_test.go",
        "resource_test.go",
        "service_test.go",
        "synonym_test.go",
        "tag_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/tag/conf:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "business.go",
        "channel.go",
        "channel_category.go",
        "channel_rule.go",
        "limit.go",
        "rank.go",
        "relation.go",
        "report.go",
        "resource.go",
        "resource_log.go",
        "service.go",
        "synonym.go",
        "tag.go",
    ],
    importpath = "go-common/app/admin/main/tag/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/tag/conf:go_default_library",
        "//app/admin/main/tag/dao:go_default_library",
        "//app/admin/main/tag/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
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
