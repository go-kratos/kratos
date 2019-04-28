package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "crm.go",
        "dao.cache.go",
        "dao.go",
        "data.go",
        "mc.cache.go",
        "memcahe.go",
        "rank.go",
        "recommend_pool.go",
        "sign.go",
    ],
    importpath = "go-common/app/interface/main/mcn/dao/mcndao",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/mcn/model:go_default_library",
        "//app/admin/main/up/util:go_default_library",
        "//app/admin/main/up/util/mathutil:go_default_library",
        "//app/interface/main/mcn/conf:go_default_library",
        "//app/interface/main/mcn/dao/cache:go_default_library",
        "//app/interface/main/mcn/dao/global:go_default_library",
        "//app/interface/main/mcn/model:go_default_library",
        "//app/interface/main/mcn/model/mcnmodel:go_default_library",
        "//app/interface/main/mcn/tool/validate:go_default_library",
        "//app/service/main/archive/api:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/pipeline/fanout:go_default_library",
        "//vendor/github.com/bluele/gcache:go_default_library",
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

go_test(
    name = "go_default_test",
    srcs = [
        "crm_test.go",
        "dao.cache_test.go",
        "dao_test.go",
        "data_test.go",
        "mc.cache_test.go",
        "memcahe_test.go",
        "rank_test.go",
        "recommend_pool_test.go",
        "sign_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/mcn/conf:go_default_library",
        "//app/interface/main/mcn/dao/global:go_default_library",
        "//app/interface/main/mcn/model:go_default_library",
        "//app/interface/main/mcn/model/mcnmodel:go_default_library",
        "//library/ecode:go_default_library",
        "//vendor/github.com/bluele/gcache:go_default_library",
        "//vendor/github.com/jinzhu/gorm:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)
