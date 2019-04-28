load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "archive_test.go",
        "audit_result_test.go",
        "dao_test.go",
        "full_test.go",
        "modules_test.go",
        "panel_test.go",
        "pgc_cond_test.go",
        "playurl_test.go",
        "region_test.go",
        "sear_inter_test.go",
        "upbfs_test.go",
        "upper_test.go",
        "user_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/tv/conf:go_default_library",
        "//app/admin/main/tv/model:go_default_library",
        "//library/log:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "archive.go",
        "audit_result.go",
        "dao.go",
        "full.go",
        "mango.go",
        "modules.go",
        "panel.go",
        "pgc_cond.go",
        "playurl.go",
        "region.go",
        "sear_inter.go",
        "upbfs.go",
        "upper.go",
        "user.go",
    ],
    importpath = "go-common/app/admin/main/tv/dao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/tv/conf:go_default_library",
        "//app/admin/main/tv/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/database/elastic:go_default_library",
        "//library/database/orm:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/time:go_default_library",
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
