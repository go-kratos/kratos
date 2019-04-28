package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "dialog_test.go",
        "jointly_test.go",
        "order_test.go",
        "panel_test.go",
        "privilege_test.go",
        "push_test.go",
        "resouce_test.go",
        "service_test.go",
        "tips_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/vip/conf:go_default_library",
        "//app/admin/main/vip/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "business.go",
        "dialog.go",
        "jointly.go",
        "month.go",
        "order.go",
        "panel.go",
        "platform.go",
        "privilege.go",
        "push.go",
        "resouce.go",
        "service.go",
        "tips.go",
        "version.go",
        "vip.go",
        "welfare.go",
    ],
    importpath = "go-common/app/admin/main/vip/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/vip/conf:go_default_library",
        "//app/admin/main/vip/dao:go_default_library",
        "//app/admin/main/vip/model:go_default_library",
        "//app/service/main/vip/model:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/google/uuid:go_default_library",
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
