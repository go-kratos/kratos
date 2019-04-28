package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "backsource_test.go",
        "dao_test.go",
        "db_operate_test.go",
        "err_msg_test.go",
        "mixed_test.go",
        "pgc_auth_test.go",
        "pgc_batch_test.go",
        "pgc_load_test.go",
        "pgc_single_test.go",
        "ugc_batch_test.go",
        "ugc_load_test.go",
        "ugc_single_test.go",
    ],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
    deps = [
        "//app/interface/main/tv/conf:go_default_library",
        "//app/interface/main/tv/model:go_default_library",
        "//library/database/sql:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "backsource.go",
        "dao.go",
        "db_operate.go",
        "err_msg.go",
        "mixed.go",
        "pgc_auth.go",
        "pgc_batch.go",
        "pgc_load.go",
        "pgc_single.go",
        "ugc_batch.go",
        "ugc_load.go",
        "ugc_single.go",
    ],
    importpath = "go-common/app/interface/main/tv/dao/cms",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/interface/main/tv/conf:go_default_library",
        "//app/interface/main/tv/model:go_default_library",
        "//library/cache/memcache:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/stat/prom:go_default_library",
        "//library/sync/errgroup:go_default_library",
        "//library/xstr:go_default_library",
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
