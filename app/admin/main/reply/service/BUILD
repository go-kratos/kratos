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
        "admin_test.go",
        "business_test.go",
        "config_test.go",
        "event_test.go",
        "moral_test.go",
        "notice_test.go",
        "reply_test.go",
        "report_test.go",
        "service_test.go",
        "subject_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/reply/conf:go_default_library",
        "//app/admin/main/reply/model:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/go-sql-driver/mysql:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "action.go",
        "admin.go",
        "business.go",
        "config.go",
        "emoji.go",
        "emoji_package.go",
        "event.go",
        "fold.go",
        "monitor.go",
        "moral.go",
        "notice.go",
        "reply.go",
        "report.go",
        "service.go",
        "subject.go",
    ],
    importpath = "go-common/app/admin/main/reply/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/reply/conf:go_default_library",
        "//app/admin/main/reply/dao:go_default_library",
        "//app/admin/main/reply/model:go_default_library",
        "//app/interface/openplatform/article/model:go_default_library",
        "//app/interface/openplatform/article/rpc/client:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/api/gorpc:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//app/service/main/relation/model:go_default_library",
        "//app/service/main/relation/rpc/client:go_default_library",
        "//app/service/main/thumbup/api:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/net/rpc/warden:go_default_library",
        "//library/queue/databus/report:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/sync/errgroup.v2:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
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
