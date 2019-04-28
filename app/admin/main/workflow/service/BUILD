package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "activity_test.go",
        "business_test.go",
        "challenge_test.go",
        "platform_test.go",
        "service_test.go",
        "wlog_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/workflow/model:go_default_library",
        "//app/admin/main/workflow/model/param:go_default_library",
        "//app/admin/main/workflow/model/search:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "activity.go",
        "business.go",
        "callback.go",
        "challenge.go",
        "event.go",
        "group.go",
        "meta.go",
        "platform.go",
        "service.go",
        "tag.go",
        "task.go",
        "wlog.go",
    ],
    importpath = "go-common/app/admin/main/workflow/service",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/credit/model/blocked:go_default_library",
        "//app/admin/main/workflow/dao:go_default_library",
        "//app/admin/main/workflow/model:go_default_library",
        "//app/admin/main/workflow/model/param:go_default_library",
        "//app/admin/main/workflow/model/search:go_default_library",
        "//app/job/main/credit/model:go_default_library",
        "//library/conf/paladin:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/golang.org/x/sync/errgroup:go_default_library",
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
