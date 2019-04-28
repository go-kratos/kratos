load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "appeal.go",
        "buss_attr.go",
        "chall.go",
        "dao.go",
        "notify.go",
        "redis.go",
        "search.go",
    ],
    importpath = "go-common/app/job/main/workflow/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/workflow/model/param:go_default_library",
        "//app/job/main/workflow/conf:go_default_library",
        "//app/job/main/workflow/model:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/elastic:go_default_library",
        "//library/database/orm:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "appeal_test.go",
        "buss_attr_test.go",
        "chall_test.go",
        "dao_test.go",
        "notify_test.go",
        "redis_test.go",
        "search_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/workflow/model/param:go_default_library",
        "//app/job/main/workflow/conf:go_default_library",
        "//app/job/main/workflow/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
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
