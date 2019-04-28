package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "apply_test.go",
        "grpc_snap_test.go",
        "grpc_test.go",
        "label_test.go",
        "mail_test.go",
        "order_admin_test.go",
        "order_test.go",
        "rank_test.go",
        "report_test.go",
        "scene_test.go",
        "script_test.go",
        "service_test.go",
        "user_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/ep/melloi/conf:go_default_library",
        "//app/admin/ep/melloi/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "apply.go",
        "bfs.go",
        "clientmoni.go",
        "cluster.go",
        "comment.go",
        "dapper.go",
        "file.go",
        "grpc.go",
        "grpc_snap.go",
        "job.go",
        "label.go",
        "mail.go",
        "order.go",
        "order_admin.go",
        "order_report.go",
        "ptest.go",
        "rank.go",
        "report.go",
        "scene.go",
        "script.go",
        "script_snap.go",
        "service.go",
        "tools.go",
        "tree.go",
        "user.go",
        "wechat.go",
    ],
    importpath = "go-common/app/admin/ep/melloi/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/ep/melloi/conf:go_default_library",
        "//app/admin/ep/melloi/dao:go_default_library",
        "//app/admin/ep/melloi/model:go_default_library",
        "//app/admin/ep/melloi/service/proto:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/robfig/cron:go_default_library",
        "//vendor/gopkg.in/gomail.v2:go_default_library",
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
        "//app/admin/ep/melloi/service/proto:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
