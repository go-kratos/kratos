package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dao.go",
        "mcn_account.go",
        "mcn_manage.go",
        "mcn_recommend.go",
        "memcache.go",
        "statistics.go",
        "up_base_Info.go",
    ],
    importpath = "go-common/app/admin/main/mcn/dao/up",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/mcn/conf:go_default_library",
        "//app/admin/main/mcn/model:go_default_library",
        "//app/interface/main/mcn/model/datamodel:go_default_library",
        "//app/interface/main/mcn/model/mcnmodel:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/net/metadata:go_default_library",
        "//library/time:go_default_library",
        "//library/xstr:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "dao_test.go",
        "mcn_account_test.go",
        "mcn_manage_test.go",
        "mcn_recommend_test.go",
        "memcache_test.go",
        "statistics_test.go",
        "up_base_Info_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/mcn/conf:go_default_library",
        "//app/admin/main/mcn/model:go_default_library",
        "//library/time:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
    ],
)
