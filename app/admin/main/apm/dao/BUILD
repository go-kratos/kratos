package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "canal_test.go",
        "dao_test.go",
        "ecode_test.go",
        "need_test.go",
        "pprof_test.go",
        "upload_test.go",
        "ut_app_test.go",
        "ut_rank_test.go",
        "ut_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/apm/conf:go_default_library",
        "//app/admin/main/apm/model/canal:go_default_library",
        "//app/admin/main/apm/model/need:go_default_library",
        "//vendor/github.com/bouk/monkey:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "canal.go",
        "dao.go",
        "ecode.go",
        "monkey.go",
        "need.go",
        "pprof.go",
        "upload.go",
        "ut.go",
        "ut_app.go",
        "ut_rank.go",
    ],
    importpath = "go-common/app/admin/main/apm/dao",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/apm/conf:go_default_library",
        "//app/admin/main/apm/model/canal:go_default_library",
        "//app/admin/main/apm/model/ecode:go_default_library",
        "//app/admin/main/apm/model/need:go_default_library",
        "//app/admin/main/apm/model/pprof:go_default_library",
        "//app/admin/main/apm/model/ut:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/conf/env:go_default_library",
        "//library/database/orm:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//vendor/github.com/bouk/monkey:go_default_library",
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
    srcs = [
        ":package-srcs",
        "//app/admin/main/apm/dao/mock:all-srcs",
    ],
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
)
