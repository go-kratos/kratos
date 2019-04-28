package(default_visibility = ["//visibility:public"])

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
        "intervs_test.go",
        "playurl_test.go",
        "region_test.go",
        "service_test.go",
        "sync_pgc_test.go",
        "uplayurl_test.go",
        "upper_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/tv/conf:go_default_library",
        "//app/admin/main/tv/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "arcType.go",
        "arc_audit.go",
        "archive.go",
        "audit_result.go",
        "full.go",
        "intervs.go",
        "label.go",
        "mango.go",
        "modules.go",
        "order.go",
        "others.go",
        "panel.go",
        "region.go",
        "sear_inter.go",
        "service.go",
        "sync_pgc.go",
        "uplayurl.go",
        "upper.go",
        "user.go",
        "video.go",
        "watermark.go",
    ],
    importpath = "go-common/app/admin/main/tv/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/tv/conf:go_default_library",
        "//app/admin/main/tv/dao:go_default_library",
        "//app/admin/main/tv/model:go_default_library",
        "//app/service/main/account/api:go_default_library",
        "//app/service/main/account/model:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//app/service/main/archive/model/archive:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
        "//vendor/github.com/pkg/errors:go_default_library",
        "//vendor/github.com/siddontang/go-mysql/mysql:go_default_library",
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
