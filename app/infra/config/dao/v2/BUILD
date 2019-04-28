package(default_visibility = ["//visibility:public"])

load(
    "@io_bazel_rules_go//go:def.bzl",
    "go_test",
    "go_library",
)

go_test(
    name = "go_default_test",
    srcs = [
        "app_test.go",
        "build_test.go",
        "config_test.go",
        "dao_test.go",
        "file_test.go",
        "force_test.go",
        "redis_test.go",
        "tag_test.go",
    ],
    embed = [":go_default_library"],
    tags = ["automanaged"],
    deps = [
        "//app/infra/config/conf:go_default_library",
        "//app/infra/config/model:go_default_library",
        "//vendor/github.com/smartystreets/goconvey/convey:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = [
        "app.go",
        "build.go",
        "config.go",
        "dao.go",
        "file.go",
        "force.go",
        "redis.go",
        "tag.go",
    ],
    importpath = "go-common/app/infra/config/dao/v2",
    tags = ["automanaged"],
    visibility = ["//visibility:public"],
    deps = [
        "//app/infra/config/conf:go_default_library",
        "//app/infra/config/model:go_default_library",
        "//library/cache/redis:go_default_library",
        "//library/database/orm:go_default_library",
        "//library/database/sql:go_default_library",
        "//library/ecode:go_default_library",
        "//library/log:go_default_library",
        "//library/xstr:go_default_library",
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
