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
        "dao_test.go",
        "mc_test.go",
        "permit_test.go",
        "reason_test.go",
        "tag_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/manager/conf:go_default_library",
        "//app/admin/main/manager/model:go_default_library",
        "//library/ecode:go_default_library",
        "//library/net/http/blademaster/middleware/permit:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "business.go",
        "dao.go",
        "journal.go",
        "mc.go",
        "permit.go",
        "reason.go",
        "tag.go",
    ],
    importpath = "go-common/app/admin/main/manager/dao",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/manager/conf:go_default_library",
        "//app/admin/main/manager/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/database/orm:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/http/blademaster/middleware/permit:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
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
