package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "api_test.go",
        "center_test.go",
        "dao_test.go",
        "field_test.go",
        "hbase_set2_test.go",
        "hbase_test.go",
        "phase_three_test.go",
        "utils_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/admin/main/up/conf:go_default_library",
        "//app/admin/main/up/model/datamodel:go_default_library",
        "//library/database/hbase.v2:go_default_library",
        "//vendor/github.com/bouk/monkey:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
        "//vendor/github.com/tsuna/gohbase/hrpc:go_default_library",
        "//vendor/gopkg.in/h2non/gock.v1:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "api.go",
        "center.go",
        "dao.go",
        "field.go",
        "hbase.go",
        "hbase_set2.go",
        "phase_three.go",
        "utils.go",
    ],
    importpath = "go-common/app/admin/main/up/dao/data",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/admin/main/up/conf:go_default_library",
        "//app/admin/main/up/model/datamodel:go_default_library",
        "//app/admin/main/up/util/hbaseutil:go_default_library",
        "//app/interface/main/creative/model/data:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/database/hbase.v2:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/net/http/blademaster:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//vendor/github.com/tsuna/gohbase/hrpc:go_default_library",
        "@org_golang_x_net//context:go_default_library",
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
