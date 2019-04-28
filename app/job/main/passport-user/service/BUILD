load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_library",
    "go_test",
)

go_library(
    name = "go_default_library",
    srcs = [
        "aes.go",
        "aso_full_migration.go",
        "aso_incr_migration.go",
        "duplicate.go",
        "passport_util.go",
        "service.go",
        "user.go",
    ],
    importpath = "go-common/app/job/main/passport-user/service",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/job/main/passport-user/conf:go_default_library",
        "//app/job/main/passport-user/dao:go_default_library",
        "//app/job/main/passport-user/model:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/log:go_default_library",
        "//library/queue/databus:go_default_library",
        "//library/queue/databus/databusutil:go_default_library",
        "//vendor/github.com/go-sql-driver/mysql:go_default_library",
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
    srcs = ["aes_test.go"],
    embed = [":go_default_library"],
    rundir = ".",
    tags = ["automanaged"],
)
