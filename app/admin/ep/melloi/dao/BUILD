package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "dao_test.go",
        "order_admin_test.go",
        "order_test.go",
        "rank_test.go",
        "scene_test.go",
        "user_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/ep/melloi/conf:go_default_library",
        "//app/admin/ep/melloi/model:go_default_library",
        "//vendor/github.com/go-sql-driver/mysql:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
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
        "dao.go",
        "dapper.go",
        "job.go",
        "label.go",
        "mail.go",
        "mysql_grpc.go",
        "mysql_grpc_snap.go",
        "order.go",
        "order_admin.go",
        "order_report.go",
        "ptestjob.go",
        "rank.go",
        "report.go",
        "reportgraph.go",
        "scene.go",
        "script.go",
        "tree.go",
        "user.go",
        "wechat.go",
    ],
    importpath = "go-common/app/admin/ep/melloi/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/ep/melloi/conf:go_default_library",
        "//app/admin/ep/melloi/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/database/orm:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/stat/prom:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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
    srcs = [":package-srcs"],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
