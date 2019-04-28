package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "bfs_test.go",
        "dao_test.go",
        "dialog_test.go",
        "jointly_test.go",
        "memcache_test.go",
        "month_test.go",
        "mysql_test.go",
        "order_test.go",
        "panel_test.go",
        "platform_test.go",
        "privilege_test.go",
        "push_test.go",
        "resouce_test.go",
        "send_test.go",
        "tips_test.go",
        "vip_test.go",
        "welfare_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/vip/conf:go_default_library",
        "//app/admin/main/vip/model:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "bfs.go",
        "dao.go",
        "dialog.go",
        "jointly.go",
        "memcache.go",
        "month.go",
        "mysql.go",
        "order.go",
        "panel.go",
        "platform.go",
        "privilege.go",
        "push.go",
        "resouce.go",
        "send.go",
        "tips.go",
        "vip.go",
        "welfare.go",
    ],
    importpath = "go-common/app/admin/main/vip/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/vip/conf:go_default_library",
        "//app/admin/main/vip/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/database/orm:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/google/uuid:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
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
